package model

import (
	"time"

	"github.com/uber/gonduit/constants"
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

// Differential represents a Phabricator Differential.
// This is a simplified version of the DifferentialRevision entity from gonduit.
type Differential struct {
	ID             string
	Title          string
	LineCount      string
	Status         constants.DifferentialStatusLegacy
	StatusName     string
	URI            string
	CreatedAt      time.Time
	ModifiedAt     time.Time
	RepositoryPHID string
}
