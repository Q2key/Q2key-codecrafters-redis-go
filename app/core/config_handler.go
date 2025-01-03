package core

func NewConfigHandler(instance Redis) *ConfigHandler {
	return &ConfigHandler{
		instance: &instance,
	}
}

type ConfigHandler struct {
	instance *Redis
}

func (h *ConfigHandler) Handle(conn Conn, args []string, _ *[]byte) {
	action, key := args[1], args[2]

	if action == "GET" && key == "dir" {
		resp := []string{key, (*h.instance).Config.Dir}
		conn.Conn().Write([]byte(StringsToRedisStrings(resp)))
		return
	}

	if action == "GET" && key == "dbfilename" {
		resp := []string{key, (*h.instance).Config.DbFileName}
		conn.Conn().Write([]byte(StringsToRedisStrings(resp)))
		return
	}

	conn.Conn().Write([]byte(ToRedisErrorString()))
}
