package action

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mevhub/internal/core/domain/game"
	"reflect"
	"time"
)

type StateChangeAction struct {
	InstanceID uuid.UUID
	State      game.State
}

func NewStateChangeAction(instanceID uuid.UUID, state game.State) *StateChangeAction {
	return &StateChangeAction{InstanceID: instanceID, State: state}
}

func (a *StateChangeAction) Perform(instance *game.LiveGameInstance) error {

	fmt.Println("Change state to", reflect.TypeOf(a.State), " at ", time.Now().UTC().String())
	instance.State = a.State
	instance.SendChange(NewStateChange(a.InstanceID, a.State))

	return nil

}
