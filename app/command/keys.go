package command

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Keys struct {
	args []string
}

func (r *Keys) FromArgs(args []string) contracts.Command[string] {
	return &Keys{
		args: args,
	}
}

func (r *Keys) Validate() bool {
	return true
}

func (r *Keys) Name() string {
	return "KEYS"
}

func (r *Keys) Args() []string {
	return r.args
}
