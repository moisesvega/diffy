package yearly

import (
	"bytes"
	"testing"
	"time"

	"github.com/moisesvega/diffy/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYearlyReporter_Report(t *testing.T) {
	// Create test data
	user1 := &entity.User{
		Username: "testuser1",
		Differentials: []*entity.Differential{
			{
				ID:        "D1",
				LineCount: 100,
				Status:    entity.Accepted,
				CreatedAt: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:        "D2",
				LineCount: 200,
				Status:    entity.Accepted,
				CreatedAt: time.Date(2023, 3, 20, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:        "D3",
				LineCount: 150,
				Status:    entity.Accepted,
				CreatedAt: time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:        "D4",
				LineCount: 300,
				Status:    entity.Accepted,
				CreatedAt: time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	user2 := &entity.User{
		Username: "testuser2",
		Differentials: []*entity.Differential{
			{
				ID:        "D5",
				LineCount: 75,
				Status:    entity.Accepted,
				CreatedAt: time.Date(2023, 6, 5, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:        "D6",
				LineCount: 125,
				Status:    entity.Accepted,
				CreatedAt: time.Date(2024, 4, 20, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	users := []*entity.User{user1, user2}

	// Create reporter and capture output
	reporter := New()
	var buf bytes.Buffer
	err := reporter.Report(users, entity.WithWriter(&buf))

	require.NoError(t, err)

	output := buf.String()

	// Check that key sections are present
	assert.Contains(t, output, "Yearly Summary")
	assert.Contains(t, output, "Year-over-Year Growth")
	assert.Contains(t, output, "Individual User Yearly Breakdown")

	// Check that both years are present
	assert.Contains(t, output, "2023")
	assert.Contains(t, output, "2024")

	// Check that usernames are present
	assert.Contains(t, output, "testuser1")
	assert.Contains(t, output, "testuser2")

	// Check that acceptance rate is not present (since we removed it)
	assert.NotContains(t, output, "Acceptance Rate")

	// Check that total row is present
	assert.Contains(t, output, "Total")

	// Check that title and emoji are not present (since we removed them)
	assert.NotContains(t, output, "Year-over-Year Analysis")
	assert.NotContains(t, output, "📊")
}

func TestYearlyReporter_ReportNoData(t *testing.T) {
	users := []*entity.User{}

	reporter := New()
	var buf bytes.Buffer
	err := reporter.Report(users, entity.WithWriter(&buf))

	require.NoError(t, err)
	assert.Empty(t, buf.String())
}

func TestYearlyReporter_ReportSingleYear(t *testing.T) {
	user := &entity.User{
		Username: "testuser",
		Differentials: []*entity.Differential{
			{
				ID:        "D1",
				LineCount: 100,
				Status:    entity.Accepted,
				CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	users := []*entity.User{user}

	reporter := New()
	var buf bytes.Buffer
	err := reporter.Report(users, entity.WithWriter(&buf))

	require.NoError(t, err)

	output := buf.String()

	// Should have yearly summary but not year-over-year growth
	assert.Contains(t, output, "Yearly Summary")

	// Should not have year-over-year growth section since there's only one year
	assert.NotContains(t, output, "Year-over-Year Growth")
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{123, "123"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
	}

	for _, test := range tests {
		result := formatNumber(test.input)
		assert.Equal(t, test.expected, result, "formatNumber(%d)", test.input)
	}
}

func TestFormatPercentage(t *testing.T) {
	tests := []struct {
		input         float64
		shouldContain string
	}{
		{25.5, "+25.5%"},
		{-10.2, "-10.2%"},
		{0.0, "0.0%"},
	}

	for _, test := range tests {
		result := formatPercentage(test.input)
		assert.Contains(t, result, test.shouldContain, "formatPercentage(%.1f)", test.input)
	}
}