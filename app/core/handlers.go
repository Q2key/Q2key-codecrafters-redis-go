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

func handleREPLCONF(h *Master, conn Conn, args []string) {
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

func handlePSYNC(h *Master, conn Conn) {
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

func handleWAIT(h *Master, conn Conn, args []string) {
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

func handleINFO(h RedisInstance, conn Conn, args []string) {
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

func handleKEYS(h RedisInstance, conn Conn, args []string) {
	keys := h.GetStore().GetKeys(args[1])
	RespondString(conn, ToRedisStrings(keys))
}

func handlePING(conn Conn) {
	RespondString(conn, ToRedisSimpleString("PONG"))
}

func handleECHO(conn Conn, args []string) {
	RespondString(conn, ToRedisSimpleString(args[1]))
}

func handleGET(h RedisInstance, conn Conn, args []string) {
	key := args[1]
	val, _ := h.GetStore().Get(key)
	if val == nil || val.IsExpired() {
		RespondString(conn, ToRedisNullBulkString())
	} else {
		RespondString(conn, ToRedisSimpleString(val.GetValue()))
	}
}

func handleCONFIG(h RedisInstance, conn Conn, args []string) {
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

func handleSET(h RedisInstance, conn Conn, args []string) {
	key, val := args[1], args[2]

	// according to the doc set is always setting the string value type
	h.GetStore().Set(key, val, STRING)

	if len(args) >= 4 {
		exp, _ := strconv.Atoi(args[4])
		h.GetStore().SetExpiredIn(key, uint64(exp))
	}

	RespondString(conn, ToRedisSimpleString("OK"))
}

func handleTYPE(h RedisInstance, conn Conn, args []string) {
	if len(args) < 1 {
		return
	}

	key := args[1]

	val, ok := h.GetStore().Get(key)
	if !ok {
		RespondString(conn, ToRedisSimpleString("none"))
	} else {
		RespondString(conn, ToRedisSimpleString(string(val.ValueType)))
	}
}

func handleXADD(h RedisInstance, conn Conn, args []string) {
	if len(args) < 1 {
		return
	}

	pairs := map[string][]string{
		"pairs": {},
	}
	for i, a := range args {
		if i == 1 {
			pairs["key"] = []string{a}
		}

		if i == 2 {
			pairs["id"] = []string{a}
		}

		if i > 2 {
			pairs["pairs"] = append(pairs["pairs"], a)
		}
	}

	key := pairs["key"]
	id := pairs["id"]

	val := &StoreValue{
		Value:     "val",
		ValueType: "stream",
	}

	h.GetStore().kvs[key[0]] = val

	RespondString(conn, ToRedisSimpleString(id[0]))
}
