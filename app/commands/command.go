package commands

type Command[T string | int] interface {
	Validate() bool
	Name() string
	Args() []T
	FromArgs(args []T) Command[T]
}
