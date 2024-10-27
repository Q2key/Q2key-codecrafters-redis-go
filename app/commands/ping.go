package commands

type Ping struct {
	args []string
}

func (r *Ping) FromArgs(args []string) Command[string] {
	return &Ping{
		args: args,
	}
}

func (r *Ping) Validate() bool {
	return len(r.args) < 2
}

func (r *Ping) Name() string {
	return "PING"
}

func (r *Ping) Args() []string {
	return r.args
}