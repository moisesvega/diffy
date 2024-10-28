package model

type Reporter interface {
	Report([]*User, ...ReporterOption) error
}
