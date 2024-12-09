package commands

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Wait struct {
	args []string
}

func (r *Wait) FromArgs(args []string) contracts.Command {
	return &Wait{
		args: args,
	}
}

func (r *Wait) Validate() bool {
	return true
}

func (r *Wait) Name() string {
	return "WAIT"
}

func (r *Wait) Args() []string {
	return r.args
}

func (r *Wait) IsWrite() bool {
	return false
}
