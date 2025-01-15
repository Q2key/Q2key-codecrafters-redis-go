package core

type CommandHandler func(instance RedisInstance, conn Conn, args []string)

func RespondString(conn Conn, data string) {
	_, err := conn.Conn().Write([]byte(data))
	if err != nil {
		// log.Fatal(err)
	}
}

func IsWriteCommand(command string) bool {
	switch command {
	case "SET":
		return true
	default:
		return false
	}
}
