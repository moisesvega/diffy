package model

type Reporter interface {
	Report([]User) error
}
