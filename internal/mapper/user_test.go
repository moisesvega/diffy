package mapper

import (
	"testing"

	"github.com/moisesvega/diffy/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/uber/gonduit/constants"
	"github.com/uber/gonduit/entities"
	"github.com/uber/gonduit/responses"
)

func Test(t *testing.T) {

}

func TestFromPhabricatorDifferential(t *testing.T) {
	tests := []struct {
		desc string
		give entities.DifferentialRevision
		want *model.Differential
	}{
		{
			desc: "success",
			give: entities.DifferentialRevision{
				ID:         "1",
				Title:      "title",
				LineCount:  "10",
				Status:     constants.DifferentialStatusLegacyAccepted,
				StatusName: "accepted",
				URI:        "uri",
			},
			want: &model.Differential{
				ID:         "1",
				Title:      "title",
				LineCount:  10,
				Status:     model.Accepted,
				StatusName: "accepted",
				URI:        "uri",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := FromPhabricatorDifferential(tt.give)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestFromPhabricatorDifferentialQueryResponse(t *testing.T) {
	tests := []struct {
		desc string
		give responses.DifferentialQueryResponse
		want []*model.Differential
	}{
		{
			desc: "success",
			give: responses.DifferentialQueryResponse{
				{
					ID:         "1",
					Title:      "title",
					LineCount:  "10",
					Status:     constants.DifferentialStatusLegacyAccepted,
					StatusName: "accepted",
					URI:        "uri",
				},
			},
			want: []*model.Differential{
				{
					ID:         "1",
					Title:      "title",
					LineCount:  10,
					Status:     model.Accepted,
					StatusName: "accepted",
					URI:        "uri",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := FromPhabricatorDifferentialQueryResponse(tt.give)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestFromPhabricatorStatus(t *testing.T) {
	tests := []struct {
		desc string
		give constants.DifferentialStatusLegacy
		want model.Status
	}{
		{
			desc: "accepted",
			give: constants.DifferentialStatusLegacyAccepted,
			want: model.Accepted,
		},
		{
			desc: "closed",
			give: constants.DifferentialStatusLegacyPublished,
			want: model.Closed,
		},
		{
			desc: "unknown",
			give: constants.DifferentialStatusLegacyNeedsReview,
			want: model.Unknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := FromPhabricatorStatus(tt.give)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestFromPhabricatorUser(t *testing.T) {
	tests := []struct {
		desc string
		give entities.User
		want *model.User
	}{
		{
			desc: "success",
			give: entities.User{
				UserName: "username",
				Email:    "email",
				PHID:     "1",
			},
			want: &model.User{
				Username: "username",
				Email:    "email",
				ID:       "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := FromPhabricatorUser(tt.give)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
