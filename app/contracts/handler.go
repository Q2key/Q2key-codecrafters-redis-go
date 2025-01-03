package contracts

type Handler interface {
	Handle(Connection, []string, *[]byte)
}
