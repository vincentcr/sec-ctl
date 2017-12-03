package sites

// Site is an interface representing the TPI client
type Site interface {
	GetID() string
	GetState() SystemState
	Exec(cmd UserCommand) error
	SubscribeToEvents() chan Event
	SubscribeToStateChange() chan StateChange
}
