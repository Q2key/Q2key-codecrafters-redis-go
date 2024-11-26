package commands

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Info struct {
	args []string
}

func (r *Info) FromArgs(args []string) contracts.Command {
	return &Info{
		args: args,
	}
}

func (r *Info) Validate() bool {
	return true
}

func (r *Info) Name() string {
	return "INFO"
}

func (r *Info) Args() []string {
	return r.args
}

func (r *Info) IsWrite() bool {
	return false
}
