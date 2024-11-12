package command

import (
	"errors"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
)

func ParseCommand(raw string) (error, *contracts.Command[string]) {
	inp := repr.ToArgs(raw)
	switch inp[0] {
	case "GET":
		cmd := new(Get).FromArgs(inp)
		return nil, &cmd
	case "SET":
		cmd := new(Set).FromArgs(inp)
		return nil, &cmd
	case "CONFIG":
		cmd := new(Config).FromArgs(inp)
		return nil, &cmd
	case "ECHO":
		cmd := new(Echo).FromArgs(inp)
		return nil, &cmd
	case "PING":
		cmd := new(Ping).FromArgs(inp)
		return nil, &cmd
	case "KEYS":
		cmd := new(Keys).FromArgs(inp)
		return nil, &cmd
	case "INFO":
		cmd := new(Info).FromArgs(inp)
		return nil, &cmd
	default:
		output := fmt.Sprintf("Unknown command: %s", inp[0])
		return errors.New(output), nil
	}
}
