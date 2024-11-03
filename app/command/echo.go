package command

type Echo struct {
	args []string
}

func (r *Echo) FromArgs(args []string) Command[string] {
	return &Echo{
		args: args,
	}
}

func (r *Echo) Validate() bool {
	return len(r.args) < 2
}

func (r *Echo) Name() string {
	return "ECHO"
}

func (r *Echo) Args() []string {
	return r.args
}
