package game

import (
	uuid "github.com/satori/go.uuid"
)

type InstanceFactory struct {
}

func (f InstanceFactory) Create(id uuid.UUID, party string, options InstanceOptions) (*Instance, error) {
	var instance = NewGameInstance()
	instance.SysID = id
	instance.PartyID = party
	instance.Options = options
	return instance, nil
}
