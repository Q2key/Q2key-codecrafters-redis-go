package command

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Config struct {
	args []string
}

func (r *Config) FromArgs(args []string) contracts.Command[string] {
	return &Config{args: args}
}

func (r *Config) Validate() bool {
	return true
}

func (r *Config) Name() string {
	return "CONFIG"
}

func (r *Config) Args() []string {
	return r.args
}
