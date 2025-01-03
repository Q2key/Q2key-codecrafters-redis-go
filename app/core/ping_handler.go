package core

func NewPingHandler(instance Redis) *PingHandler {
	return &PingHandler{
		instance: instance,
	}
}

type PingHandler struct {
	instance Redis
}

func (h *PingHandler) Handle(conn Conn, _ []string, _ *[]byte) {
	conn.Conn().Write([]byte(FromStringToRedisCommonString("PONG")))
}
