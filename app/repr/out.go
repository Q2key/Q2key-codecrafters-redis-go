package repr

import "fmt"

func FromStringsArray(input []string) string {
	n := len(input)
	r := fmt.Sprintf("*%d", n)
	for i := 0; i < n; i++ {
		s := input[i]
		r += fmt.Sprintf("\r\n$%d\r\n%s", len(s), s)
	}
	r += "\r\n"
	return r
}

func ErrorString() string {
	return "$-1\r\n"
}

func ErrorStringWithMessage(error error) string {
	return fmt.Sprintf("$-1%s\r\n", error.Error())
}

func BulkString(input string) string {
	l := len(input)
	return fmt.Sprintf("$%d\r\n%s\r\n", l, input)
}

func FromString(input string) string {
	return fmt.Sprintf("+%s\r\n", input)
}
