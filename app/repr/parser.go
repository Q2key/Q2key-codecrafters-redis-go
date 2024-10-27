package repr

import (
	"errors"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"strconv"
)

const (
	NA = iota
	Error
	String
	BulkStrings
	Number
	Array
)

func ParseType(q string) int {
	fb := q[0:1]
	// error
	switch fb {
	case "*":
		return Array
	case ":":
		return Number
	case "$":
		return BulkStrings
	default:
		return NA
	}
}

func ParseArray(q string) []string {
	s := q[1:]
	sli := make([]string, 0)
	for i := 0; i < len(s); i++ {
		if string(s[i]) == "$" {
			ni := i + 1
			nj := i + 2

			sl, err := strconv.Atoi(s[ni:nj])
			if sl == 0 || err != nil {
				break
			}

			st := nj + 2
			fi := st + sl

			if fi > len(s) {
				break
			}

			sli = append(sli, s[st:fi])
		}
	}

	return sli
}

func ParseCommand(raw string) (error, *commands.Command[string]) {
	inp := ParseArray(raw)
	switch inp[0] {
	case "GET":
		cmd := new(commands.Get).FromArgs(inp)
		return nil, &cmd
	case "SET":
		cmd := new(commands.Set).FromArgs(inp)
		return nil, &cmd
	case "CONFIG":
		cmd := new(commands.Config).FromArgs(inp)
		return nil, &cmd
	case "ECHO":
		cmd := new(commands.Echo).FromArgs(inp)
		return nil, &cmd
	case "PING":
		cmd := new(commands.Ping).FromArgs(inp)
		return nil, &cmd
	case "KEYS":
		cmd := new(commands.Keys).FromArgs(inp)
		return nil, &cmd
	default:
		return errors.New(fmt.Sprintf("Unknown command: %s", inp[0])), nil
	}
}

func ToRegularString(input string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", input))
}

func ToStringArray(input []string) []byte {
	n := len(input)
	r := fmt.Sprintf("*%d", n)
	for i := 0; i < n; i++ {
		s := input[i]
		r += fmt.Sprintf("\r\n$%d\r\n%s", len(s), s)
	}
	r += "\r\n"
	return []byte(r)
}

func ToErrorString(input *string) []byte {
	if input != nil {
		return make([]byte, 0)
	}
	return []byte("$-1\r\n")
}
