package heatmap

import (
	"bytes"
	"testing"
	"time"

	"github.com/moisesvega/diffy/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	r := New()
	require.NotNil(t, r)
}

func TestReport(t *testing.T) {
	today := time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC)
	give := []*model.User{
		{
			Username: "username",
			Differentials: []*model.Differential{
				{
					Title:      "title",
					URI:        "uri",
					LineCount:  11,
					ModifiedAt: today.AddDate(0, 0, -1), // yesterday
				},
			},
		},
	}

	want := `╭─────────┬───╮
│         │Oct│
├─────────┼───┤
│ Sunday  │ 1 │
│ Monday  │ 0 │
│ Tuesday │   │
│Wednesday│   │
│Thursday │   │
│ Friday  │   │
│Saturday │   │
╰─────────┴───╯
`
	r := &reporter{
		now: func() time.Time {
			return today
		},
	}
	w := &bytes.Buffer{}
	err := r.Report(give,
		model.WithWriter(w),
		model.WithSince(
			// 7 days ago
			time.Now().AddDate(0, 0, -7),
		),
	)
	require.NoError(t, err)
	assert.NotEmpty(t, w)
	assert.Equal(t, want, w.String())

}
