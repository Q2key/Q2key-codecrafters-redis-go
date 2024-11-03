package command

type Keys struct {
	args []string
}

func (r *Keys) FromArgs(args []string) Command[string] {
	return &Keys{
		args: args,
	}
}

func (r *Keys) Validate() bool {
	return len(r.args) <= 2
}

func (r *Keys) Name() string {
	return "KEYS"
}

func (r *Keys) Args() []string {
	return r.args
}
