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

// IsActive tells whether the machine is in an "active" state. Active is
// defined in legacy vagrant as "having an id file in the data dir" which is
// roughtly equivalent to "has an entry in the underlying provider."
func (s State) IsActive() bool {
	switch s {
	case CREATED, HALTED, PENDING:
		return true
	case UNKNOWN, DESTROYED, NOT_CREATED:
		return false
	default:
		return false
	}

}
