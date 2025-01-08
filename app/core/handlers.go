package core

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func RespondString(conn Conn, data string) {
	_, err := conn.Conn().Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}
}

const (
	MockDbContent = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"
)

func handleReplconf(r RedisInstance, conn Conn, args []string) {
	master, ok := r.(*Master)
	if !ok {
		return
	}
	h := master
	if len(args) > 2 && args[1] == "ACK" {
		cnt := args[2]
		num, _ := strconv.Atoi(cnt)
		id := conn.Id()

		h.UpdateReplica(id, num)
		*h.AckChan <- Ack{ConnId: id, Offset: num}
	} else {
		RespondString(conn, ToRedisSimpleString("OK"))
	}
}

func handlePsync(r RedisInstance, conn Conn, _ []string) {
	master, ok := r.(*Master)
	if !ok {
		return
	}
	h := master

	h.RegisterReplicaConn(conn)

	resp := ToRedisSimpleString(fmt.Sprintf("FULLRESYNC %s 0", conn.Id()))
	rdbBuff, err := hex.DecodeString(MockDbContent)
	if err != nil {
		log.Fatal(err)
	}

	var sb strings.Builder
	defer sb.Reset()
	sb.WriteString(resp)
	sb.WriteString("$88\r\n")
	sb.Write(rdbBuff)

	RespondString(conn, sb.String())
}

func handleWaitAsMaster(r RedisInstance, conn Conn, args []string) {
	master, ok := r.(*Master)
	if !ok {
		return
	}

	h := master
	rep, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}

	timeout, err := strconv.Atoi(args[2])
	if err != nil {
		return
	}

	done := map[string]bool{}

	for _, r := range *h.RepConnPool {
		if r.Offset() >= *h.ReceivedBytes {
			done[r.Id()] = true
			continue
		}

		go RespondString(r, "*3\r\n$8\r\nREPLCONF\r\n$6\r\nGETACK\r\n$1\r\n*\r\n")

	}

	ch := *h.AckChan

awaitingLoop:
	for len(done) < rep {
		select {
		case c := <-ch:
			{
				if c.Offset >= *h.ReceivedBytes {
					done[c.ConnId] = true
				}
			}
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			{
				break awaitingLoop
			}
		}
	}

	RespondString(conn, ToRedisInteger(strconv.Itoa(len(done))))
}

func handleInfo(h RedisInstance, conn Conn, args []string) {
	r := h.GetConfig().GetReplica()
	res := "role:master"
	if r != nil {
		res = "role:slave"
	}

	for _, v := range args {
		if v == "replication" {
			res += ":master_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
			res += ":master_repl_offset:0"
		}
	}

	RespondString(conn, ToRedisBulkString(res))
}

func handleKeys(h RedisInstance, conn Conn, args []string) {
	keys := h.GetStore().GetKeys(args[1])
	RespondString(conn, ToRedisStrings(keys))
}

func handleGet(h RedisInstance, conn Conn, args []string) {
	key := args[1]
	val, _ := h.GetStore().Get(key)
	if val == nil || val.IsExpired() {
		RespondString(conn, ToRedisNullBulkString())
	} else {
		RespondString(conn, ToRedisSimpleString(val.ToString()))
	}
}

func handleConfig(h RedisInstance, conn Conn, args []string) {
	action, key := args[1], args[2]

	if action == "GET" && key == "dir" {
		resp := []string{key, h.GetConfig().Dir}
		RespondString(conn, ToRedisStrings(resp))
		return
	}

	if action == "GET" && key == "dbfilename" {
		resp := []string{key, h.GetConfig().DbFileName}
		RespondString(conn, ToRedisStrings(resp))
		return
	}

	RespondString(conn, ToRedisNullBulkString())
}

func handleSet(h RedisInstance, conn Conn, args []string) {
	key, val := args[1], args[2]

	// according to the doc set is always setting the string value type
	h.GetStore().Set(key, val, STRING)

	if len(args) >= 4 {
		exp, _ := strconv.Atoi(args[4])
		h.GetStore().SetExpiredIn(key, uint64(exp))
	}

	RespondString(conn, ToRedisSimpleString("OK"))
}

func handleType(h RedisInstance, conn Conn, args []string) {
	if len(args) < 1 {
		return
	}

	key := args[1]

	val, ok := h.GetStore().Get(key)
	if !ok {
		RespondString(conn, ToRedisSimpleString("none"))
	} else {
		RespondString(conn, ToRedisSimpleString(string(val.GetType())))
	}
}

func handleXadd(h RedisInstance, conn Conn, args []string) {
	if len(args) < 3 {
		return
	}

	parts := strings.Split(args[2], "-")

	var err error
	var msTime float64 = -1
	seqNum := -1

	store := h.GetStore()
	storeKey := args[1]
	payload := strings.Join(args[3:], ":")

	if args[2] == "*" {
		msTime = float64(time.Now().UnixMilli())
		seqNum = 0
		value := NewStreamValue(msTime)

		value.WriteSequence(msTime, seqNum, payload)
		store.SetRedisValue(storeKey, value)

		RespondString(conn, ToRedisBulkString(value.ToString()))
		return
	}

	msTime, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return
	}

	if parts[1] != "*" {
		seqNum, err = strconv.Atoi(parts[1])
		if err != nil {
			return
		}
	}

	stream, ok := store.Get(storeKey)
	if !ok {
		value := NewStreamValue(msTime)
		if seqNum == -1 {
			seqNum = value.UpdateSeqKey(msTime)
		}

		value.WriteSequence(msTime, seqNum, payload)
		store.SetRedisValue(storeKey, value)

		RespondString(conn, ToRedisSimpleString(value.ToString()))
		return
	}

	value, _ := stream.(*StreamValue)
	if !value.KeyExists(msTime) {
		value.WriteSequence(msTime, seqNum, payload)
	}

	if seqNum == -1 {
		seqNum = value.UpdateSeqKey(msTime)
	}

	canSave, cause := value.CanSave(msTime, seqNum)
	if canSave {

		value.WriteSequence(msTime, seqNum, payload)
		store.SetRedisValue(storeKey, value)

		RespondString(conn, ToRedisSimpleString(value.ToString()))

	} else {
		RespondString(conn, ToSimpleError(*cause))
	}
}

func handlePing(_ RedisInstance, conn Conn, _ []string) {
	RespondString(conn, ToRedisSimpleString("PONG"))
}

func handleEcho(_ RedisInstance, conn Conn, args []string) {
	RespondString(conn, ToRedisSimpleString(args[1]))
}

func handleXrange(_ RedisInstance, conn Conn, args []string) {
	fmt.Println(args)
	RespondString(conn, ToRedisSimpleString("HELLO"))
}
