package phabricator

import (
	"github.com/uber/gonduit/constants"
	"github.com/uber/gonduit/entities"
)

type User struct {
	entities.User
	Differentials          []*entities.DifferentialRevision
	Reviews                []*entities.DifferentialRevision
	publishedDifferentials []*entities.DifferentialRevision
}

func (u *User) PublishedDifferentials() []*entities.DifferentialRevision {
	if u.publishedDifferentials != nil {
		return u.publishedDifferentials
	}

	out := make([]*entities.DifferentialRevision, 0)
	for _, d := range u.Differentials {
		if d.Status == constants.DifferentialStatusLegacyPublished {
			out = append(out, d)
		}
	}
	u.publishedDifferentials = out
	return out
}
