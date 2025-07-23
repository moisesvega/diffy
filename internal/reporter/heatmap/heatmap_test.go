package heatmap

import (
	"bytes"
	"errors"
	"fmt"
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
Total Differentials: 16 | Total Lines Changed: 176
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

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{123, "123"},
		{1000, "1,000"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
		{12345678, "12,345,678"},
		{123456789, "123,456,789"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("input_%d", test.input), func(t *testing.T) {
			result := formatNumber(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestReportWithLargeNumbers(t *testing.T) {
	// Test the formatted totals output with larger numbers
	wednesday := time.Date(2024, 10, 30, 0, 0, 0, 0, time.UTC)

	// Create many differentials to get larger totals
	var differentials []*entity.Differential
	for i := 0; i < 5000; i++ { // 5,000 differentials
		differentials = append(differentials, &entity.Differential{
			Title:      fmt.Sprintf("title_%d", i),
			URI:        fmt.Sprintf("uri_%d", i),
			LineCount:  250,                             // 250 lines each = 1,250,000 total lines
			ModifiedAt: wednesday.AddDate(0, 0, -i%365), // Spread over a year
		})
	}

	give := []*entity.User{
		{
			Username:      "test_user",
			Differentials: differentials,
		},
	}

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
		entity.WithSince(wednesday.AddDate(-1, 0, 0)), // 1 year ago
	)
	require.NoError(t, err)

	output := w.String()
	// Check that the output contains formatted numbers with commas
	assert.Contains(t, output, "Total Differentials: 5,000")
	assert.Contains(t, output, "Total Lines Changed: 1,250,000")
}
