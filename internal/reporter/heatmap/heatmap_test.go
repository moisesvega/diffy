package heatmap

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/moisesvega/diffy/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	require.NotNil(t, New())
}

func TestReport(t *testing.T) {
	// Given a user with differentials spread across 4 days:
	// - Wednesday (today): 7 differentials (3*2 + 1)
	// - Tuesday: 5 differentials (2*2 + 1)
	// - Monday: 3 differentials (1*2 + 1)
	// - Sunday: 1 differential (0*2 + 1)
	// And we're viewing the heatmap on Wednesday (today)
	// When we generate the heatmap
	// Then we should see:
	// - A table starting from Sunday
	// - Activity counts matching the differentials per day
	// - Proper month headers
	// - Empty cells for days without activity (Thursday, Friday, Saturday)
	wednesday := time.Date(2024, 10, 30, 0, 0, 0, 0, time.UTC)
	tuesday := wednesday.AddDate(0, 0, -1)
	monday := wednesday.AddDate(0, 0, -2)
	sunday := wednesday.AddDate(0, 0, -3)

	collection := []time.Time{wednesday, tuesday, monday, sunday}

	differentials := make([]*entity.Differential, 0)
	for i, day := range collection {
		for j := 0; j <= i*2; j++ {
			differentials = append(differentials, &entity.Differential{
				Title:      "title",
				URI:        "uri",
				LineCount:  11,
				ModifiedAt: day,
			})
		}

	}
	give := []*entity.User{
		{
			Username:      "username",
			Differentials: differentials,
		},
	}

	want := `
╭─────────┬───┬───╮
│username │Oct│   │
├─────────┼───┼───┤
│ Sunday  │ 0 │ 7 │
│ Monday  │ 0 │ 5 │
│ Tuesday │ 0 │ 3 │
│Wednesday│ 0 │ 1 │
│Thursday │ 0 │   │
│ Friday  │ 0 │   │
│Saturday │ 0 │   │
╰─────────┴───┴───╯
`
	originalNow := now
	t.Cleanup(func() {
		now = originalNow
	})
	now = func() time.Time {
		return wednesday
	}
	r := &reporter{}
	w := &bytes.Buffer{}
	err := r.Report(give,
		entity.WithWriter(w),
		entity.WithSince(
			// 7 days ago
			wednesday.AddDate(0, 0, -7),
		),
	)
	require.NoError(t, err)
	assert.NotEmpty(t, w)
	assert.Equal(t, want, w.String())
}

type failedWriter struct {
	err error
}

func (f failedWriter) Write(p []byte) (n int, err error) {
	return 0, f.err
}

func TestReporterError(t *testing.T) {
	today := time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC)
	originalNow := now
	t.Cleanup(func() {
		now = originalNow
	})
	now = func() time.Time {
		return today
	}
	r := &reporter{}
	give := []*entity.User{
		{
			Username: "username",
			Differentials: []*entity.Differential{
				{
					Title:      "title",
					URI:        "uri",
					LineCount:  11,
					ModifiedAt: today.AddDate(0, 0, -1), // yesterday
				},
			},
		},
	}

	want := errors.New("sad")
	err := r.Report(give, entity.WithWriter(failedWriter{err: want}))
	require.Error(t, err)
	assert.ErrorIs(t, err, want)
}
