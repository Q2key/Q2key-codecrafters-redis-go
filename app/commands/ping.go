package commands

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Ping struct {
	args []string
}

func (r *Ping) FromArgs(args []string) contracts.Command {
	return &Ping{
		args: args,
	}
}

func (r *Ping) Validate() bool {
	return true
}

func (r *Ping) Name() string {
	return "PING"
}

func (r *Ping) Args() []string {
	return r.args
}
