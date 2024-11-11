package command

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Echo struct {
	args []string
}

func (r *Echo) FromArgs(args []string) contracts.Command[string] {
	return &Echo{
		args: args,
	}
}

func (r *Echo) Validate() bool {
	return true
}

func (r *Echo) Name() string {
	return "ECHO"
}

func (r *Echo) Args() []string {
	return r.args
}
