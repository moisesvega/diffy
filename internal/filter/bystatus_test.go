package filter

import (
	"slices"
	"testing"

	"github.com/moisesvega/diffy/internal/model"
	"github.com/stretchr/testify/require"
)

func TestByStatus(t *testing.T) {

	tests := []struct {
		desc   string
		status model.Status
		give   []*model.Differential
		want   []*model.Differential
	}{
		{
			desc:   "success",
			status: model.Closed,
			give: []*model.Differential{
				{
					ID:         "1",
					Title:      "title",
					LineCount:  10,
					Status:     model.Closed,
					StatusName: "closed",
					URI:        "uri",
				},
				{
					ID:         "2",
					Title:      "title",
					LineCount:  10,
					Status:     model.Accepted,
					StatusName: "accepted",
					URI:        "uri",
				},
			},
			want: []*model.Differential{
				{
					ID:         "1",
					Title:      "title",
					LineCount:  10,
					Status:     model.Closed,
					StatusName: "closed",
					URI:        "uri",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			out := slices.DeleteFunc(tt.give, ByStatus(tt.status))
			require.Len(t, out, len(tt.want))
		})
	}
}
