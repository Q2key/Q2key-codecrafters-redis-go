package commands

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Psync struct {
	args []string
}

func (r *Psync) FromArgs(args []string) contracts.Command {
	return &Psync{
		args: args,
	}
}

func (r *Psync) Validate() bool {
	return true
}

func (r *Psync) Name() string {
	return "PSYNC"
}

func (r *Psync) Args() []string {
	return r.args
}
