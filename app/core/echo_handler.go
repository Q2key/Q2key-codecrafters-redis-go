package core

func NewEchoHandler(instance Redis) *EchoHandler {
	return &EchoHandler{
		instance: instance,
	}
}

type EchoHandler struct {
	instance Redis
}

func (h *EchoHandler) Handle(conn Conn, args []string, _ *[]byte) {
	conn.Conn().Write([]byte(FromStringToRedisCommonString(args[1])))
}
