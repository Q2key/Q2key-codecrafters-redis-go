package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var excludedMarks = map[string]bool{
	"*": true,
	"-": true,
	"+": true,
}

func isIdPair(mark string) bool {
	return excludedMarks[mark] || strings.Contains(mark, "-")
}

func tryGetBlockParams(args []string) (bool, *int) {
	if args[1] != "block" {
		return false, nil
	}

	var val int
	var err error

	val, err = strconv.Atoi(args[2])
	if err != nil {
		return false, nil
	}

	return true, &val
}

func getArgMap(args []string) [][]string {
	piv := 0
	tmp := args[2:]
	for i := 1; i < len(tmp); i++ {
		if isIdPair(tmp[i]) && !isIdPair(tmp[i-1]) {
			piv = i
		}
	}

	keys := tmp[:piv]
	vals := tmp[piv:]

	n := len(keys)
	if n != len(vals) {
		return [][]string{}
	}

	out := make([][]string, n)
	for i := 0; i < n; i++ {
		out[i] = []string{keys[i], vals[i]}
	}

	return out
}

func parseId(from string) (float64, int) {
	var ts float64
	var id int
	if excludedMarks[from] {
		ts = -1
		id = -1
	} else {
		parts := strings.Split(from, "-")
		ts, _ = strconv.ParseFloat(parts[0], 64)
		id, _ = strconv.Atoi(parts[1])
	}
	return ts, id
}

func handleXread(ins RedisInstance, conn Conn, args []string) {
	needBlock, blockForMs := tryGetBlockParams(args)

	_ = needBlock
	_ = blockForMs

	if needBlock {
		key := args[4]
		tsi, seq := parseId(args[5])

		s, _ := ins.GetStore().Get(key)
		v, _ := s.(*StreamValue)

		_ = tsi
		_ = seq

		go func() {
		loop:
			for {
				select {
				case <-*v.C:
					fmt.Println("sending updated value (new value occured)")

					/* res */
					var sb strings.Builder
					sb.WriteString(ToArrayDefString(1))
					buildResponse(&sb, v, key, tsi)
					RespondString(conn, sb.String())

					break loop
				case <-time.After(time.Duration(*blockForMs) * time.Millisecond):
					fmt.Println("sending existing value (time is over)")

					/* res */
					var sb strings.Builder
					sb.WriteString(ToArrayDefString(1))
					buildResponse(&sb, v, key, tsi)
					RespondString(conn, sb.String())

					break loop
				}
			}

			RespondString(conn, "$-1\r\n")
		}()

		return
	}

	argMap := getArgMap(args)

	/* */
	var sb strings.Builder
	sb.WriteString(ToArrayDefString(len(argMap)))
	for _, a := range argMap {
		key, val := a[0], a[1]
		tsi, _ := parseId(val)
		v, _ := ins.GetStore().Get(key)
		rv, _ := v.(*StreamValue)

		buildResponse(&sb, rv, key, tsi)
	}

	RespondString(conn, sb.String())
}

func buildResponse(sb *strings.Builder, rv *StreamValue, key string, ts float64) {
	k := formKey(ts, rv.LastSidx)
	val := rv.Paris[k]

	sb.WriteString(ToArrayDefString(2))
	sb.WriteString(ToRedisBulkString(key))
	sb.WriteString(ToArrayDefString(1))
	sb.WriteString(ToArrayDefString(2))
	sb.WriteString(ToRedisBulkString(k))
	sb.WriteString(ToRedisStrings(val))
}
