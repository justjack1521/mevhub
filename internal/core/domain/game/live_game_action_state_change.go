package game

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"time"
)

type StateChangeAction struct {
	InstanceID uuid.UUID
	State      State
}

func NewStateChangeAction(instanceID uuid.UUID, state State) *StateChangeAction {
	return &StateChangeAction{InstanceID: instanceID, State: state}
}

func (a *StateChangeAction) Perform(game *LiveGameInstance) error {

	fmt.Println("Change state to", reflect.TypeOf(a.State), " at ", time.Now().UTC().String())
	game.State = a.State

	game.SendChange(NewStateChange(a.InstanceID, a.State))

	return nil

}
