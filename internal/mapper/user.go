package mapper

import (
	"strconv"
	"time"

	"github.com/moisesvega/diffy/internal/entity"
	"github.com/uber/gonduit/constants"
	"github.com/uber/gonduit/entities"
	"github.com/uber/gonduit/responses"
)

// FromPhabricatorUser maps an entities.User to a entity.User.
func FromPhabricatorUser(in entities.User) *entity.User {
	return &entity.User{
		Username: in.UserName,
		Email:    in.Email,
		ID:       in.PHID,
	}
}

// FromPhabricatorDifferential maps an entities.DifferentialRevision to a entity.Differential.
func FromPhabricatorDifferential(in entities.DifferentialRevision) *entity.Differential {
	count, err := strconv.Atoi(in.LineCount)
	if err != nil {
		count = 0
	}
	return &entity.Differential{
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

// FromPhabricatorStatus maps a constants.DifferentialStatusLegacy to a entity.Status.
func FromPhabricatorStatus(in constants.DifferentialStatusLegacy) entity.Status {
	switch in {
	case constants.DifferentialStatusLegacyAccepted:
		return entity.Accepted
	case constants.DifferentialStatusLegacyPublished:
		return entity.Closed
	default:
		return entity.Unknown
	}
}

// FromPhabricatorDifferentialQueryResponse maps an entities.User to a entity.User.
func FromPhabricatorDifferentialQueryResponse(in responses.DifferentialQueryResponse) []*entity.Differential {
	out := make([]*entity.Differential, 0, len(in))
	for _, d := range in {
		out = append(out, FromPhabricatorDifferential(*d))
	}
	return out
}
