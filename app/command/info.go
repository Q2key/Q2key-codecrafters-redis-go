package command

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Info struct {
	args []string
}

func (r *Info) FromArgs(args []string) contracts.Command[string] {
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
