package mapper

import (
	"github.com/moisesvega/diffy/internal/model"
	"github.com/uber/gonduit/entities"
	"github.com/uber/gonduit/responses"
)

// FromPhabricatorUser maps an entities.User to a model.User.
func FromPhabricatorUser(in entities.User) *model.User {
	return &model.User{
		Username: in.UserName,
		Email:    in.Email,
		ID:       in.PHID,
	}
}

// FromPhabricatorDifferential maps an entities.DifferentialRevision to a model.Differential.
func FromPhabricatorDifferential(in entities.DifferentialRevision) *model.Differential {
	return &model.Differential{
		ID:        in.PHID,
		Title:     in.Title,
		LineCount: in.LineCount,
	}
}

// FromPhabricatorDifferentialQueryResponse maps an entities.User to a model.User.
func FromPhabricatorDifferentialQueryResponse(in responses.DifferentialQueryResponse) []*model.Differential {
	out := make([]*model.Differential, 0, len(in))
	for _, d := range in {
		out = append(out, FromPhabricatorDifferential(*d))
	}
	return out
}
