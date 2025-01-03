package core

func NewGetHandler(instance Redis) *GetHandler {
	return &GetHandler{
		instance: instance,
	}
}

type GetHandler struct {
	instance Redis
}

func (h *GetHandler) Handle(conn Conn, args []string, _ *[]byte) {
	key := args[1]
	val, _ := h.instance.Store.Get(key)
	if val == nil || val.IsExpired() {
		conn.Conn().Write([]byte(ToRedisErrorString()))
	} else {
		conn.Conn().Write([]byte(FromStringToRedisCommonString(val.GetValue())))
	}
}
