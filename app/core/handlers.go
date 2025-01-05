package core

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

func handleINFO(h *Redis, conn Conn, args []string) {
	r := h.Config.GetReplica()
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

	conn.Conn().Write([]byte(FromStringToRedisBulkString(res)))
}

func handleKEYS(h *Redis, conn Conn, args []string) {
	keys := h.Store.GetKeys(args[1])
	conn.Conn().Write([]byte(StringsToRedisStrings(keys)))
}

func handlePING(conn Conn) {
	conn.Conn().Write([]byte(FromStringToRedisCommonString("PONG")))
}

func handleECHO(conn Conn, args []string) {
	conn.Conn().Write([]byte(FromStringToRedisCommonString(args[1])))
}

func handleGET(h *Redis, conn Conn, args []string) {
	key := args[1]
	val, _ := h.Store.Get(key)
	if val == nil || val.IsExpired() {
		conn.Conn().Write([]byte(ToRedisErrorString()))
	} else {
		conn.Conn().Write([]byte(FromStringToRedisCommonString(val.GetValue())))
	}
}

func handlePSYNC(h *Redis, conn Conn) {
	h.RegisterReplicaConn(conn)
	mess := fmt.Sprintf("FULLRESYNC %s 0", conn.Id())
	resp := FromStringToRedisCommonString(mess)
	rdb := "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"
	rdbBuff, _ := hex.DecodeString(rdb)
	chunkA := []byte(resp)
	chunkB := []byte("$88\r\n")
	chunkC := rdbBuff
	res := CombineBuffers(chunkA, chunkB, chunkC)
	conn.Conn().Write(res)
}

func handleCONFIG(h *Redis, conn Conn, args []string) {
	action, key := args[1], args[2]

	if action == "GET" && key == "dir" {
		resp := []string{key, (*h).Config.Dir}
		conn.Conn().Write([]byte(StringsToRedisStrings(resp)))
		return
	}

	if action == "GET" && key == "dbfilename" {
		resp := []string{key, (*h).Config.DbFileName}
		conn.Conn().Write([]byte(StringsToRedisStrings(resp)))
		return
	}

	conn.Conn().Write([]byte(ToRedisErrorString()))
}

func handleREPLCONF(h *Redis, conn Conn, args []string) {
	if len(args) > 2 && args[1] == "ACK" {
		cnt := args[2]
		num, _ := strconv.Atoi(cnt)
		id := conn.Id()

		h.UpdateReplica(id, num)
		*h.AckChan <- Ack{ConnId: id, Offset: num}

		conn.Conn().Write([]byte(""))
	} else {
		conn.Conn().Write([]byte(FromStringToRedisCommonString("OK")))
	}
}

func handleSET(h *Redis, conn Conn, args []string, bytes *[]byte) {
	key, val := args[1], args[2]

	valueTypes := GetValueTypes(string(*bytes))
	valueType, ok := valueTypes[key]
	if !ok {
		return
	}

	h.Store.Set(key, val, valueType)

	if len(args) >= 4 {
		exp, _ := strconv.Atoi(args[4])
		h.Store.SetExpiredIn(key, uint64(exp))
	}

	conn.Conn().Write([]byte(FromStringToRedisCommonString("OK")))
}

func handleTYPE(h *Redis, conn Conn, args []string) {
	if len(args) < 1 {
		return
	}

	key := args[1]

	val, ok := h.Store.Get(key)
	if !ok {
		conn.Conn().Write([]byte(FromStringToRedisCommonString("none")))
	} else {
		conn.Conn().Write([]byte(FromStringToRedisCommonString(string(val.ValueType))))
	}
}

func handleWAIT(h *Redis, conn Conn, args []string) {
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

		go func() {
			r.Conn().Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$6\r\nGETACK\r\n$1\r\n*\r\n"))
		}()

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

	v := strconv.Itoa(len(done))
	conn.Conn().Write([]byte(FromStringToRedisInteger(v)))
}

func handleXADD(h *Redis, conn Conn, args []string) {
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

	h.Store.kvs[key[0]] = val

	conn.Conn().Write([]byte(FromStringToRedisCommonString(id[0])))
}
