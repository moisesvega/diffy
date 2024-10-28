package model

import (
	"time"
)

// User represents a Phabricator User.
// This is a simplified version of the User entity from gonduit.
type User struct {
	Username      string
	Email         string
	ID            string
	Differentials []*Differential
	Reviews       []*Differential
}

type Status string

const (
	Accepted Status = "accepted"
	Closed   Status = "closed"
	Unknown  Status = "unknown"
)

// Differential represents a Phabricator Differential.
// This is a simplified version of the DifferentialRevision entity from gonduit.
type Differential struct {
	ID             string
	Title          string
	LineCount      int
	Status         Status
	StatusName     string
	URI            string
	CreatedAt      time.Time
	ModifiedAt     time.Time
	RepositoryPHID string
}
