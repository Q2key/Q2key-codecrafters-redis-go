package commands

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Type struct {
	args []string
}

func (r *Type) FromArgs(args []string) contracts.Command {
	return &Type{
		args: args,
	}
}

func (r *Type) Validate() bool {
	return true
}

func (r *Type) Name() string {
	return "PING"
}

func (r *Type) Args() []string {
	return r.args
}

func (r *Type) IsWrite() bool {
	return false
}
