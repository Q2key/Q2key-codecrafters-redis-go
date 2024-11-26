package commands

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Get struct {
	args []string
}

func (r *Get) FromArgs(c []string) contracts.Command {
	return &Get{
		args: c,
	}
}

func (r *Get) Validate() bool {
	return true
}

func (r *Get) Name() string {
	return "GET"
}

func (r *Get) Args() []string {
	return r.args
}

func (r *Get) IsWrite() bool {
	return false
}
