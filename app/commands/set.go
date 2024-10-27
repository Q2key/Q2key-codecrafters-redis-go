package commands

type Set struct {
	args []string
}

func (r *Set) FromArgs(args []string) Command[string] {
	return &Set{
		args: args,
	}
}

func (r *Set) Validate() bool {
	return len(r.args) < 2
}

func (r *Set) Name() string {
	return "SET"
}

func (r *Set) Args() []string {
	return r.args
}