package model

type Client interface {
	// GetUsers returns a list of users with their differentials and reviews.
	GetUsers(strings []string) ([]User, error)
}
