package server

type ErrorHandlerDefault struct {
}

func NewErrorHandlerDefault() *ErrorHandlerDefault {
	return &ErrorHandlerDefault{}
}

func (d *ErrorHandlerDefault) Handle(svr *GameServer, err error) {
	svr.errorCount++
}
