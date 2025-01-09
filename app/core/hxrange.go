package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func handleXrange(ins RedisInstance, conn Conn, args []string) {
	key, from, to := args[1], args[2], args[3]
	fromParts := strings.Split(from, "-")
	fromTs, err := strconv.ParseFloat(fromParts[0], 64)
	if err != nil {
		return
	}

	fromSeq, err := strconv.Atoi(fromParts[1])
	if err != nil {
		return
	}

	toParts := strings.Split(to, "-")
	toTs, err := strconv.ParseFloat(toParts[0], 64)
	if err != nil {
		return
	}

	toSeq, err := strconv.Atoi(toParts[1])
	if err != nil {
		return
	}

	res, err := xrange(ins, key, fromTs, toTs, fromSeq, toSeq)
	if err != nil {
		RespondString(conn, ToSimpleError(err.Error()))
	} else {
		RespondString(conn, res)
	}
}

func xrange(
	ins RedisInstance,
	key string,
	fromTs, toTs float64,
	fromSeq, toSeq int,
) (string, error) {
	v, ok := ins.GetStore().Get(key)
	if !ok {
		return "", errors.New("something went wrong")
	}

	rv, ok := v.(*StreamValue)
	if !ok {
		return "", errors.New("something went wrong")
	}

	keys := []string{}
	for ts, v := range rv.Value {
		if ts >= fromTs && ts <= toTs {
			for _, idx := range v {
				if idx >= fromSeq && idx <= toSeq {
					keys = append(keys, formKey(ts, idx))
				}
			}
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", len(keys)))
	for _, k := range keys {
		values := rv.Paris[k]
		sb.WriteString(fmt.Sprintf("*%d\r\n", len(values)))
		sb.WriteString(ToRedisBulkString(k))
		sb.WriteString(ToRedisStrings(values))
	}

	sb.WriteString("\r\n")
	return sb.String(), nil
}
