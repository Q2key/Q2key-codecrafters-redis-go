package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func handleXread(ins RedisInstance, conn Conn, args []string) {
	if args[1] != "streams" {
		return
	}
	key, from := args[2], args[3]

	fmt.Println(args)
	var fromTs float64
	var err error
	var fromSeq int
	var res *string

	fromParts := strings.Split(from, "-")
	fromTs, err = strconv.ParseFloat(fromParts[0], 64)
	if err != nil {
		return
	}

	res, err = xread(ins, key, fromTs, fromSeq)
	if err != nil {
		RespondString(conn, ToSimpleError(err.Error()))
	} else {
		RespondString(conn, *res)
	}
}

func xread(
	ins RedisInstance,
	key string,
	fromTs float64, fromSeq int,
) (*string, error) {
	v, ok := ins.GetStore().Get(key)
	if !ok {
		return nil, errors.New("something went wrong")
	}

	rv, ok := v.(*StreamValue)
	if !ok {
		return nil, errors.New("something went wrong")
	}

	k := formKey(fromTs, rv.LastSidx)
	val := (*rv).Paris[k]

	return nil, nil
}
