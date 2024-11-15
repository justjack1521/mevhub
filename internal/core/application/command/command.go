package command

import "github.com/justjack1521/mevium/pkg/mevent"

type BasicCommand struct {
	events []mevent.Event
}

func (x *BasicCommand) GetQueuedEvents() []mevent.Event {
	return x.events
}

func (x *BasicCommand) QueueEvent(evt mevent.Event) {
	x.events = append(x.events, evt)
}
