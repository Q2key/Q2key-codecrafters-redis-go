package redis

type Handler struct {
	Store Store
}

func NewHandler() Handler {
	return Handler{
		Store: NewRedisStore(),
	}
}
