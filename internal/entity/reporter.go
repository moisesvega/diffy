package entity

type Reporter interface {
	Report([]*User, ...ReporterOption) error
}
