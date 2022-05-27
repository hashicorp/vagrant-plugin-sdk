package core

type State uint

const (
	UNKNOWN State = iota
	CREATED
	DESTROYED
	HALTED
	NOT_CREATED
	PENDING
)
