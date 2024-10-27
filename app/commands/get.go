package commands

type Get struct {
	args []string
}

func (r *Get) FromArgs(c []string) Command[string] {
	return &Get{
		args: c,
	}
}

func (r *Get) Validate() bool {
	if len(r.args) < 2 {
		return false
	}

	return false
}

func (r *Get) Name() string {
	return "GET"
}

func (r *Get) Args() []string {
	return r.args
}
