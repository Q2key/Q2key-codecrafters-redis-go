package commands

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Set struct {
	args []string
}

func (r *Set) FromArgs(args []string) contracts.Command {
	return &Set{
		args: args,
	}
}

func (r *Set) Validate() bool {
	return true
}

func (r *Set) Name() string {
	return "SET"
}

func (r *Set) Args() []string {
	return r.args
}
