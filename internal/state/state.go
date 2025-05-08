package state

// All possible client states
type State uint8

const (
	Null State = iota
	Status
	Login
	Configuration
	Play
)
