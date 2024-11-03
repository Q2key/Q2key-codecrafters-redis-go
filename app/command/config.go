package command

type Config struct {
	args []string
}

func (r *Config) FromArgs(args []string) Command[string] {
	return &Config{args: args}
}

func (r *Config) Validate() bool {
	return len(r.args) < 3
}

func (r *Config) Name() string {
	return "CONFIG"
}

func (r *Config) Args() []string {
	return r.args
}
