package handler

type IHandler interface {
	Login()
}

type HandlerImpl struct {
}

func InitHandler() *IHandler {
	return
}
