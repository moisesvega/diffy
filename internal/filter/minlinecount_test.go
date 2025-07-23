package filter

import (
	"slices"
	"testing"

	"github.com/moisesvega/diffy/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestByLineCount(t *testing.T) {
	tests := []struct {
		desc  string
		count int
		give  []*entity.Differential
		want  []*entity.Differential
	}{
		{
			desc:  "success",
			count: 10,
			give: []*entity.Differential{
				{
					ID:         "1",
					Title:      "title",
					LineCount:  9,
					Status:     entity.Closed,
					StatusName: "closed",
					URI:        "uri",
				},
				{
					ID:         "2",
					Title:      "title",
					LineCount:  20,
					Status:     entity.Accepted,
					StatusName: "accepted",
					URI:        "uri",
				},
			},
			want: []*entity.Differential{
				{
					ID:         "2",
					Title:      "title",
					LineCount:  20,
					Status:     entity.Accepted,
					StatusName: "accepted",
					URI:        "uri",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := slices.DeleteFunc(tt.give, MinLineCount(tt.count))
			require.Len(t, got, len(tt.want))
			assert.EqualValues(t, tt.want, got)
		})
	}
}
