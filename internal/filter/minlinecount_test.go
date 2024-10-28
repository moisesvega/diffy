package filter

import (
	"slices"
	"testing"

	"github.com/moisesvega/diffy/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestByLineCount(t *testing.T) {

	tests := []struct {
		desc  string
		count int
		give  []*model.Differential
		want  []*model.Differential
	}{
		{
			desc:  "success",
			count: 10,
			give: []*model.Differential{
				{
					ID:         "1",
					Title:      "title",
					LineCount:  9,
					Status:     model.Closed,
					StatusName: "closed",
					URI:        "uri",
				},
				{
					ID:         "2",
					Title:      "title",
					LineCount:  20,
					Status:     model.Accepted,
					StatusName: "accepted",
					URI:        "uri",
				},
			},
			want: []*model.Differential{
				{
					ID:         "2",
					Title:      "title",
					LineCount:  20,
					Status:     model.Accepted,
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
