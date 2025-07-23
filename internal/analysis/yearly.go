package analysis

import (
	"time"

	"github.com/moisesvega/diffy/internal/entity"
)

// YearlyStats represents statistics for a specific year
type YearlyStats struct {
	Year                  int     `json:"year"`
	TotalDifferentials    int     `json:"total_differentials"`
	TotalLinesChanged     int     `json:"total_lines_changed"`
	AvgLinesPerDiff       float64 `json:"avg_lines_per_diff"`
	AcceptedDifferentials int     `json:"accepted_differentials"`
	AcceptanceRate        float64 `json:"acceptance_rate"`
}

// YearOverYearProgress represents year-over-year progression data
type YearOverYearProgress struct {
	CurrentYear          YearlyStats `json:"current_year"`
	PreviousYear         YearlyStats `json:"previous_year"`
	DifferentialsGrowth  float64     `json:"differentials_growth_percent"`
	LinesChangedGrowth   float64     `json:"lines_changed_growth_percent"`
	AcceptanceRateChange float64     `json:"acceptance_rate_change_percent"`
}

// AnalyzeTotalDifferentialPerYear analyzes differential statistics grouped by year
func AnalyzeTotalDifferentialPerYear(users []*entity.User) map[int]YearlyStats {
	yearlyData := make(map[int]YearlyStats)

	for _, user := range users {
		for _, diff := range user.Differentials {
			year := diff.CreatedAt.Year()

			stats := yearlyData[year]
			stats.Year = year
			stats.TotalDifferentials++
			stats.TotalLinesChanged += diff.LineCount

			if diff.Status == entity.Accepted {
				stats.AcceptedDifferentials++
			}

			yearlyData[year] = stats
		}
	}

	// Calculate derived metrics
	for year, stats := range yearlyData {
		if stats.TotalDifferentials > 0 {
			stats.AvgLinesPerDiff = float64(stats.TotalLinesChanged) / float64(stats.TotalDifferentials)
			stats.AcceptanceRate = float64(stats.AcceptedDifferentials) / float64(stats.TotalDifferentials) * 100
		}
		yearlyData[year] = stats
	}

	return yearlyData
}

// CalculateYearOverYearProgress calculates year-over-year progression for a specific year
func CalculateYearOverYearProgress(yearlyData map[int]YearlyStats, targetYear int) *YearOverYearProgress {
	currentStats, currentExists := yearlyData[targetYear]
	previousStats, previousExists := yearlyData[targetYear-1]

	if !currentExists || !previousExists {
		return nil
	}

	progress := &YearOverYearProgress{
		CurrentYear:  currentStats,
		PreviousYear: previousStats,
	}

	// Calculate growth percentages
	if previousStats.TotalDifferentials > 0 {
		progress.DifferentialsGrowth = float64(currentStats.TotalDifferentials-previousStats.TotalDifferentials) /
			float64(previousStats.TotalDifferentials) * 100
	}

	if previousStats.TotalLinesChanged > 0 {
		progress.LinesChangedGrowth = float64(currentStats.TotalLinesChanged-previousStats.TotalLinesChanged) /
			float64(previousStats.TotalLinesChanged) * 100
	}

	progress.AcceptanceRateChange = currentStats.AcceptanceRate - previousStats.AcceptanceRate

	return progress
}

// GetAvailableYears returns all years that have data, sorted in descending order
func GetAvailableYears(yearlyData map[int]YearlyStats) []int {
	years := make([]int, 0, len(yearlyData))
	for year := range yearlyData {
		years = append(years, year)
	}

	// Sort years in descending order (most recent first)
	for i := 0; i < len(years)-1; i++ {
		for j := i + 1; j < len(years); j++ {
			if years[i] < years[j] {
				years[i], years[j] = years[j], years[i]
			}
		}
	}

	return years
}

// AnalyzeUserYearlyProgress analyzes year-over-year progress for a specific user
func AnalyzeUserYearlyProgress(user *entity.User) map[int]YearlyStats {
	return AnalyzeTotalDifferentialPerYear([]*entity.User{user})
}

// GetCurrentYearStats returns statistics for the current year
func GetCurrentYearStats(users []*entity.User) YearlyStats {
	currentYear := time.Now().Year()
	yearlyData := AnalyzeTotalDifferentialPerYear(users)

	if stats, exists := yearlyData[currentYear]; exists {
		return stats
	}

	return YearlyStats{Year: currentYear}
}
