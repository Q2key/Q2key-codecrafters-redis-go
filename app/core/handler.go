package core

type Handler interface {
	Handle(Conn, []string, *[]byte)
}
