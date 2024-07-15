package app

type GameApplication struct {
	Queries     *GameApplicationQueries
	Commands    *GameApplicationCommands
	consumers   []ApplicationConsumer
	subscribers []ApplicationSubscriber
}

type GameApplicationQueries struct {
}

type GameApplicationCommands struct {
}

type GameApplicationTranslators struct {
}

func NewGameApplication(core *CoreApplication) *GameApplication {
	var application = &GameApplication{
		consumers:   []ApplicationConsumer{},
		subscribers: []ApplicationSubscriber{},
	}
	application.Queries = &GameApplicationQueries{}
	application.Commands = &GameApplicationCommands{}
	return application
}
