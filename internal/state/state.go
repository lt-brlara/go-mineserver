package state

// All possible client states
const (
	StateNull SessionState = iota
	StateStatus
	StateLogin
	StateConfiguration
	StatePlay
)
