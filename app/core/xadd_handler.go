package core

import "fmt"

func NewXaddHandler(instance Redis) *XaddHandler {
	return &XaddHandler{
		instance: instance,
	}
}

type XaddHandler struct {
	instance Redis
}

func (h *XaddHandler) Handle(conn Conn, args []string, _ *[]byte) {
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

	fmt.Println(pairs)

	key := pairs["key"]
	id := pairs["id"]

	val := &StoreValue{
		Value:     "val",
		ValueType: "stream",
	}

	h.instance.Store.kvs[key[0]] = val

	conn.Conn().Write([]byte(FromStringToRedisCommonString(id[0])))
}
