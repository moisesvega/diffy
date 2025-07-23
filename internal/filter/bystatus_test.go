package filter

import (
	"slices"
	"testing"

	"github.com/moisesvega/diffy/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestByStatus(t *testing.T) {
	tests := []struct {
		desc   string
		status entity.Status
		give   []*entity.Differential
		want   []*entity.Differential
	}{
		{
			desc:   "success",
			status: entity.Closed,
			give: []*entity.Differential{
				{
					ID:         "1",
					Title:      "title",
					LineCount:  10,
					Status:     entity.Closed,
					StatusName: "closed",
					URI:        "uri",
				},
				{
					ID:         "2",
					Title:      "title",
					LineCount:  10,
					Status:     entity.Accepted,
					StatusName: "accepted",
					URI:        "uri",
				},
			},
			want: []*entity.Differential{
				{
					ID:         "1",
					Title:      "title",
					LineCount:  10,
					Status:     entity.Closed,
					StatusName: "closed",
					URI:        "uri",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := slices.DeleteFunc(tt.give, ByStatus(tt.status))
			require.Len(t, got, len(tt.want))
			assert.EqualValues(t, tt.want, got)
		})
	}
}
