package ports

import (
	"mevhub/internal/app"
)

type MultiGrpcServerImplementation struct {
	app *app.Application
}

func NewMultiGrpcServerImplementation(app *app.Application) *MultiGrpcServerImplementation {
	return &MultiGrpcServerImplementation{app: app}
}
