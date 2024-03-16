package app

import "github.com/justjack1521/mevium/pkg/mevent"

type ApplicationSubscriber interface {
	mevent.Handler
}
