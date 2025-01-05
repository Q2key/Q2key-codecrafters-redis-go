package core

func IsWriteCommand(command string) bool {
	switch command {
	case "SET":
		return true
	default:
		return false
	}
}
