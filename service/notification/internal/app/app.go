package app

type (
	EventHandler interface {
		CreateSubscription() error
		DeleteSubscription() error
		SendNotification() error
	}
	app struct {
		eventHandler EventHandler
	}
)

func NewApp(eventHandler EventHandler) *app {
	return &app{
		eventHandler: eventHandler,
	}
}

func (a *app) Run() error {

	return nil
}
