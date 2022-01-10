package core

type Push interface {
	Push() (err error)
}
