package mapper

import (
	"strconv"
	"time"

	"github.com/moisesvega/diffy/internal/model"
	"github.com/uber/gonduit/constants"
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
	count, err := strconv.Atoi(in.LineCount)
	if err != nil {
		count = 0
	}
	return &model.Differential{
		ID:         in.ID,
		Title:      in.Title,
		LineCount:  count,
		Status:     FromPhabricatorStatus(in.Status),
		StatusName: in.StatusName,
		URI:        in.URI,
		CreatedAt:  time.Time(in.DateCreated),
		ModifiedAt: time.Time(in.DateModified),
	}
}

// FromPhabricatorStatus maps a constants.DifferentialStatusLegacy to a model.Status.
func FromPhabricatorStatus(in constants.DifferentialStatusLegacy) model.Status {
	switch in {
	case constants.DifferentialStatusLegacyAccepted:
		return model.Accepted
	case constants.DifferentialStatusLegacyPublished:
		return model.Closed
	default:
		return model.Unknown
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
