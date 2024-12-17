package contracts

type Handler interface {
	Handle(Connector, []string, *[]byte)
}
