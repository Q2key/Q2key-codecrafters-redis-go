package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func handleXrange(ins RedisInstance, conn Conn, args []string) {
	key, from, to := args[1], args[2], args[3]

	toParts := strings.Split(to, "-")

	var fromTs float64
	var toTs float64
	var err error
	var fromSeq int
	var toSeq int
	var res *string

	if from == "-" {
		fromTs = -1
		fromSeq = -1
	} else {
		fromParts := strings.Split(from, "-")
		fromTs, err = strconv.ParseFloat(fromParts[0], 64)
		if err != nil {
			return
		}

		fromSeq, err = strconv.Atoi(fromParts[1])
		if err != nil {
			return
		}
	}

	if to == "+" {
		toTs = -1
		toSeq = -1
	} else {
		toTs, err = strconv.ParseFloat(toParts[0], 64)
		if err != nil {
			return
		}

		toSeq, err = strconv.Atoi(toParts[1])
		if err != nil {
			return
		}
	}

	res, err = xrange(ins, key, fromTs, toTs, fromSeq, toSeq)
	if err != nil {
		RespondString(conn, ToSimpleError(err.Error()))
	} else {
		RespondString(conn, *res)
	}
}

func xrange(
	ins RedisInstance,
	key string,
	fromTs, toTs float64,
	fromSeq, toSeq int,
) (*string, error) {
	v, ok := ins.GetStore().Get(key)
	if !ok {
		return nil, errors.New("something went wrong")
	}

	rv, ok := v.(*StreamValue)
	if !ok {
		return nil, errors.New("something went wrong")
	}

	var keys []string
	if fromTs == -1 {
		keys = xrangeTill(*rv, toTs, toSeq)
	} else if toTs == -1 {
		keys = xrangeFrom(*rv, fromTs, fromSeq)
	} else {
		keys = xrangeRage(*rv, fromTs, toTs, fromSeq, toSeq)
	}

	result := write(*rv, keys)
	return &result, nil
}

func xrangeRage(
	rv StreamValue,
	fromTs, toTs float64,
	fromSeq, toSeq int,
) []string {
	var keys []string
	for ts, v := range rv.Value {
		if ts >= fromTs && ts <= toTs {
			for _, idx := range v {
				if idx >= fromSeq && idx <= toSeq {
					keys = append(keys, formKey(ts, idx))
				}
			}
		}
	}
	return keys
}

func xrangeTill(
	rv StreamValue, toTs float64, toSeq int,
) []string {
	var keys []string
	for ts, v := range rv.Value {
		if ts <= toTs {
			for _, idx := range v {
				if idx <= toSeq {
					keys = append(keys, formKey(ts, idx))
				}
			}
		}
	}

	return keys
}

func xrangeFrom(
	rv StreamValue, fromTs float64, fromSeq int,
) []string {
	var keys []string
	for ts, v := range rv.Value {
		if ts >= fromTs {
			for _, idx := range v {
				if idx >= fromSeq {
					keys = append(keys, formKey(ts, idx))
				}
			}
		}
	}

	return keys
}

func write(rv StreamValue, keys []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", len(keys)))
	for _, k := range keys {
		values := rv.Paris[k]
		sb.WriteString(fmt.Sprintf("*%d\r\n", len(values)))
		sb.WriteString(ToRedisBulkString(k))
		sb.WriteString(ToRedisStrings(values))
	}

	sb.WriteString("\r\n")
	return sb.String()
}
