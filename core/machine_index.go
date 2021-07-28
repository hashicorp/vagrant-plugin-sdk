package core

type TargetIndex interface {
	Delete(machine Machine) (err error)
	Get(uuid string) (entry Machine, err error)
	Includes(uuid string) (exists bool, err error)
	Set(entry Machine) (updatedEntry Machine, err error)
	Recover(entry Machine) (updatedEntry Machine, err error)
}
