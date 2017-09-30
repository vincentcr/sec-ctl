package site

// Client is an interface representing the TPI client
type Client interface {
	Exec(cmd UserCommand) error
	SubscribeToEvents() chan Event
	SubscribeToStateChange() chan StateChange
	GetState() SystemState
}
