package command

type Get struct {
	args []string
}

func (r *Get) FromArgs(c []string) Command[string] {
	return &Get{
		args: c,
	}
}

func (r *Get) Validate() bool {
	return true
}

func (r *Get) Name() string {
	return "GET"
}

func (r *Get) Args() []string {
	return r.args
}
