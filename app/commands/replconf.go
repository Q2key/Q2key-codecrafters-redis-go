package commands

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type ReplConf struct {
	args []string
}

func (r *ReplConf) FromArgs(args []string) contracts.Command[string] {
	return &ReplConf{
		args: args,
	}
}

func (r *ReplConf) Validate() bool {
	return true
}

func (r *ReplConf) Name() string {
	return "REPLCONF"
}

func (r *ReplConf) Args() []string {
	return r.args
}