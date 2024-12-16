package commands

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type TypeCmd struct {
	args []string
}

func (r *TypeCmd) FromArgs(args []string) contracts.Command {
	return &TypeCmd{
		args: args,
	}
}

func (r *TypeCmd) Validate() bool {
	return true
}

func (r *TypeCmd) Name() string {
	return "TYPE"
}

func (r *TypeCmd) Args() []string {
	return r.args
}

func (r *TypeCmd) IsWrite() bool {
	return false
}
