package core

func NewKeysHandler(instance Redis) *KeysHandler {
	return &KeysHandler{
		instance: instance,
	}
}

type KeysHandler struct {
	instance Redis
}

func (h *KeysHandler) Handle(conn Conn, args []string, _ *[]byte) {
	t := args[1]
	keys := h.instance.Store.GetKeys(t)

	conn.Conn().Write([]byte(StringsToRedisStrings(keys)))
}
