package core

import (
	"errors"
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

	res, _ := xrange(ins, key, fromTs, toTs, fromSeq, toSeq)
	RespondString(conn, res)
}

func xrange(ins RedisInstance, key string, fromTs, toTs float64, fromSeq, toSeq int) (string, error) {
	v, ok := ins.GetStore().Get(key)
	if !ok {
		return "", errors.New("something went wrong")
	}

	rv, ok := v.(*StreamValue)
	if !ok {
		return "", errors.New("something went wrong")
	}

	/*
		[
		  [
		    "1526985054069-0",
		    [
		      "temperature",
		      "36",
		      "humidity",
		      "95"
		    ]
		  ],
		  [
		    "1526985054079-0",
		    [
		      "temperature",
		      "37",
		      "humidity",
		      "94"
		    ]
		  ],
		]
	*/

	for ts, v := range rv.Value {
		if ts >= fromTs && ts <= toTs {
			tmp := []string{}
			_ = tmp
			for i, idx := range v {
				_, _ = i, idx
			}
		}
	}

	return "not implemented yet", nil
}
