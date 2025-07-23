package analysis

import (
	"testing"
	"time"

	"github.com/moisesvega/diffy/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeTotalDifferentialPerYear(t *testing.T) {
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
				Status:    entity.Closed,
				CreatedAt: time.Date(2023, 3, 20, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:        "D3",
				LineCount: 150,
				Status:    entity.Accepted,
				CreatedAt: time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	user2 := &entity.User{
		Username: "testuser2",
		Differentials: []*entity.Differential{
			{
				ID:        "D4",
				LineCount: 75,
				Status:    entity.Accepted,
				CreatedAt: time.Date(2023, 6, 5, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:        "D5",
				LineCount: 300,
				Status:    entity.Accepted,
				CreatedAt: time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	users := []*entity.User{user1, user2}

	// Test the function
	result := AnalyzeTotalDifferentialPerYear(users)

	// Verify 2023 data
	stats2023, exists := result[2023]
	require.True(t, exists, "Expected data for year 2023")

	assert.Equal(t, 3, stats2023.TotalDifferentials)
	assert.Equal(t, 375, stats2023.TotalLinesChanged)
	assert.Equal(t, 2, stats2023.AcceptedDifferentials)

	expectedAcceptanceRate := float64(2) / float64(3) * 100
	assert.Equal(t, expectedAcceptanceRate, stats2023.AcceptanceRate)

	expectedAvgLines := float64(375) / float64(3)
	assert.Equal(t, expectedAvgLines, stats2023.AvgLinesPerDiff)

	// Verify 2024 data
	stats2024, exists := result[2024]
	require.True(t, exists, "Expected data for year 2024")

	assert.Equal(t, 2, stats2024.TotalDifferentials)
	assert.Equal(t, 450, stats2024.TotalLinesChanged)
	assert.Equal(t, 2, stats2024.AcceptedDifferentials)
	assert.Equal(t, 100.0, stats2024.AcceptanceRate)
}

func TestCalculateYearOverYearProgress(t *testing.T) {
	yearlyData := map[int]YearlyStats{
		2023: {
			Year:                 2023,
			TotalDifferentials:   10,
			TotalLinesChanged:    1000,
			AcceptedDifferentials: 8,
			AcceptanceRate:       80.0,
			AvgLinesPerDiff:      100.0,
		},
		2024: {
			Year:                 2024,
			TotalDifferentials:   15,
			TotalLinesChanged:    1200,
			AcceptedDifferentials: 14,
			AcceptanceRate:       93.33,
			AvgLinesPerDiff:      80.0,
		},
	}

	progress := CalculateYearOverYearProgress(yearlyData, 2024)
	require.NotNil(t, progress, "Expected progress data")

	// Test differentials growth: (15-10)/10 * 100 = 50%
	assert.Equal(t, 50.0, progress.DifferentialsGrowth)

	// Test lines changed growth: (1200-1000)/1000 * 100 = 20%
	assert.Equal(t, 20.0, progress.LinesChangedGrowth)

	// Test acceptance rate change: 93.33 - 80.0 = 13.33%
	expectedAcceptanceChange := 13.33
	assert.InDelta(t, expectedAcceptanceChange, progress.AcceptanceRateChange, 0.01)
}

func TestCalculateYearOverYearProgressMissingData(t *testing.T) {
	yearlyData := map[int]YearlyStats{
		2024: {
			Year:               2024,
			TotalDifferentials: 15,
			TotalLinesChanged:  1200,
		},
	}

	// Test with missing previous year data
	progress := CalculateYearOverYearProgress(yearlyData, 2024)
	assert.Nil(t, progress, "Expected nil progress when previous year data is missing")

	// Test with missing current year data
	progress = CalculateYearOverYearProgress(yearlyData, 2025)
	assert.Nil(t, progress, "Expected nil progress when current year data is missing")
}

func TestGetAvailableYears(t *testing.T) {
	yearlyData := map[int]YearlyStats{
		2022: {Year: 2022},
		2024: {Year: 2024},
		2023: {Year: 2023},
		2021: {Year: 2021},
	}

	years := GetAvailableYears(yearlyData)
	expected := []int{2024, 2023, 2022, 2021}

	assert.Equal(t, expected, years)
}

func TestAnalyzeUserYearlyProgress(t *testing.T) {
	user := &entity.User{
		Username: "testuser",
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
				Status:    entity.Closed,
				CreatedAt: time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	result := AnalyzeUserYearlyProgress(user)

	assert.Len(t, result, 2)

	stats2023, exists := result[2023]
	require.True(t, exists)
	assert.Equal(t, 1, stats2023.TotalDifferentials)

	stats2024, exists := result[2024]
	require.True(t, exists)
	assert.Equal(t, 1, stats2024.TotalDifferentials)
}

func TestGetCurrentYearStats(t *testing.T) {
	currentYear := time.Now().Year()
	user := &entity.User{
		Username: "testuser",
		Differentials: []*entity.Differential{
			{
				ID:        "D1",
				LineCount: 100,
				Status:    entity.Accepted,
				CreatedAt: time.Now(),
			},
		},
	}

	stats := GetCurrentYearStats([]*entity.User{user})

	assert.Equal(t, currentYear, stats.Year)
	assert.Equal(t, 1, stats.TotalDifferentials)
}