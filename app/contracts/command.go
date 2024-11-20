package contracts

type Command interface {
	Validate() bool
	Name() string
	Args() []string
	FromArgs(args []string) Command
}
