package core

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

	key := args[1]

	conn.Conn().Write([]byte(FromStringToRedisCommonString(key)))
}
