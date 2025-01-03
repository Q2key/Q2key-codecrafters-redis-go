package core

func NewTypeHandler(instance Redis) *TypeHandler {
	return &TypeHandler{
		instance: instance,
	}
}

type TypeHandler struct {
	instance Redis
}

func (h *TypeHandler) Handle(conn Conn, args []string, _ *[]byte) {
	if len(args) < 1 {
		return
	}

	key := args[1]

	val, ok := h.instance.Store.Get(key)
	if !ok {
		conn.Conn().Write([]byte(FromStringToRedisCommonString("none")))
	} else {
		conn.Conn().Write([]byte(FromStringToRedisCommonString(val.ValueType)))
	}
}
