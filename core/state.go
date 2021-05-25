package core

type State uint

const (
	UNKNOWN State = iota
	PENDING
	CREATED
	DESTROYED
)
