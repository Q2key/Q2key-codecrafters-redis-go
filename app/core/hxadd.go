package core

import (
	"strconv"
	"strings"
	"time"
)

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
		*value.C <- true
	} else {
		RespondString(conn, ToSimpleError(*cause))
	}
}
