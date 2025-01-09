package core

import (
	"errors"
	"strconv"
	"strings"
)

func handleXread(ins RedisInstance, conn Conn, args []string) {
	key, from := args[1], args[2]

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

	keys := []string{formKey(fromTs, fromSeq)}

	result := write(*rv, keys)
	return &result, nil
}
