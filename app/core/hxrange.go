package core

import (
	"errors"
	"strings"
)

func handleXrange(ins RedisInstance, conn Conn, args []string) {
	key, iidx, jidx := args[1], args[2], args[3]
	fts, fs := parseId(iidx)
	tts, ts := parseId(jidx)

	res, err := xrange(ins, key, fts, tts, fs, ts)
	if err != nil {
		RespondString(conn, ToSimpleError(err.Error()))
	} else {
		RespondString(conn, *res)
	}
}

func xrange(
	ins RedisInstance,
	key string,
	its, iseq float64,
	jts, jset int,
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
	if its == -1 {
		keys = xrangeTill(*rv, iseq, jset)
	} else if iseq == -1 {
		keys = xrangeFrom(*rv, its, jts)
	} else {
		keys = xrangeRage(*rv, its, iseq, jts, jset)
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
	sb.WriteString(ToArrayDefString(len(keys)))
	for _, k := range keys {
		values := rv.Paris[k]
		sb.WriteString(ToArrayDefString(len(values)))
		sb.WriteString(ToRedisBulkString(k))
		sb.WriteString(ToRedisStrings(values))
	}

	sb.WriteString(CLRF)
	return sb.String()
}
