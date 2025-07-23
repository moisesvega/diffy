package yearly

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/moisesvega/diffy/internal/analysis"
	"github.com/moisesvega/diffy/internal/entity"
)

const (
	_background     = "#212830"
	_whiteFontColor = "#FAFAFA"
	_greenColor     = "#26a641"
	_redColor       = "#f85149"
	_blueColor      = "#58a6ff"
	_yellowColor    = "#f9e71e"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(_whiteFontColor)).
			Background(lipgloss.Color(_background)).
			Bold(true).
			Padding(1, 0)

	positiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(_greenColor)).
			Bold(true)

	negativeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(_redColor)).
			Bold(true)

	neutralStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(_blueColor)).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(_yellowColor)).
			Bold(true)
)

type reporter struct{}

// formatNumber formats an integer with commas for thousands separators
func formatNumber(n int) string {
	str := strconv.Itoa(n)
	if len(str) <= 3 {
		return str
	}

	var result strings.Builder
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}
	return result.String()
}

// formatPercentage formats a percentage with appropriate color coding
func formatPercentage(value float64) string {
	percentStr := fmt.Sprintf("%.1f%%", value)
	
	if value > 0 {
		return positiveStyle.Render("+" + percentStr)
	} else if value < 0 {
		return negativeStyle.Render(percentStr)
	}
	return neutralStyle.Render(percentStr)
}

func (r *reporter) Report(users []*entity.User, option ...entity.ReporterOption) error {
	opts := &entity.ReporterOptions{}
	for _, o := range option {
		o(opts)
	}

	var w io.Writer = os.Stdout
	if opts.Writer != nil {
		w = opts.Writer
	}

	// Analyze yearly data for all users combined
	yearlyData := analysis.AnalyzeTotalDifferentialPerYear(users)
	if len(yearlyData) == 0 {
		return nil
	}

	// Get available years
	years := analysis.GetAvailableYears(yearlyData)
	
	// Create yearly summary table
	fmt.Fprint(w, "\n")
	r.printYearlySummary(w, yearlyData, years)

	// Print year-over-year progressions if we have multiple years
	if len(years) > 1 {
		fmt.Fprint(w, "\n")
		r.printYearOverYearProgress(w, yearlyData, years)
	}

	// Print individual user yearly breakdowns
	if len(users) > 1 {
		fmt.Fprint(w, "\n")
		r.printUserYearlyBreakdowns(w, users)
	}

	return nil
}

func (r *reporter) printYearlySummary(w io.Writer, yearlyData map[int]analysis.YearlyStats, years []int) {
	fmt.Fprint(w, headerStyle.Render("Yearly Summary"), "\n")

	headers := []string{"Year", "Differentials", "Lines Changed", "Avg Lines/Diff"}
	rows := make([][]string, 0, len(years)+1)

	// Calculate totals
	totalDifferentials := 0
	totalLinesChanged := 0

	for _, year := range years {
		stats := yearlyData[year]
		totalDifferentials += stats.TotalDifferentials
		totalLinesChanged += stats.TotalLinesChanged
		
		rows = append(rows, []string{
			strconv.Itoa(year),
			formatNumber(stats.TotalDifferentials),
			formatNumber(stats.TotalLinesChanged),
			fmt.Sprintf("%.1f", stats.AvgLinesPerDiff),
		})
	}

	// Add total row
	avgLinesOverall := float64(0)
	if totalDifferentials > 0 {
		avgLinesOverall = float64(totalLinesChanged) / float64(totalDifferentials)
	}
	
	rows = append(rows, []string{
		headerStyle.Render("Total"),
		headerStyle.Render(formatNumber(totalDifferentials)),
		headerStyle.Render(formatNumber(totalLinesChanged)),
		headerStyle.Render(fmt.Sprintf("%.1f", avgLinesOverall)),
	})

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Background(lipgloss.Color(_background))).
		Headers(headers...).
		Rows(rows...)

	fmt.Fprint(w, t.Render(), "\n")
}

func (r *reporter) printYearOverYearProgress(w io.Writer, yearlyData map[int]analysis.YearlyStats, years []int) {
	fmt.Fprint(w, headerStyle.Render("Year-over-Year Growth"), "\n")

	headers := []string{"Period", "Differentials Growth", "Lines Growth"}
	rows := make([][]string, 0)

	for i := 0; i < len(years)-1; i++ {
		currentYear := years[i]
		previousYear := years[i+1]

		progress := analysis.CalculateYearOverYearProgress(yearlyData, currentYear)
		if progress != nil {
			period := fmt.Sprintf("%d vs %d", currentYear, previousYear)
			rows = append(rows, []string{
				period,
				formatPercentage(progress.DifferentialsGrowth),
				formatPercentage(progress.LinesChangedGrowth),
			})
		}
	}

	if len(rows) > 0 {
		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Background(lipgloss.Color(_background))).
			Headers(headers...).
			Rows(rows...)

		fmt.Fprint(w, t.Render(), "\n")
	}
}

func (r *reporter) printUserYearlyBreakdowns(w io.Writer, users []*entity.User) {
	fmt.Fprint(w, headerStyle.Render("Individual User Yearly Breakdown"), "\n")

	for _, user := range users {
		userYearlyData := analysis.AnalyzeUserYearlyProgress(user)
		if len(userYearlyData) == 0 {
			continue
		}

		fmt.Fprintf(w, "\n%s:\n", neutralStyle.Render(user.Username))

		years := analysis.GetAvailableYears(userYearlyData)
		headers := []string{"Year", "Differentials", "Lines Changed"}
		rows := make([][]string, 0, len(years))

		for _, year := range years {
			stats := userYearlyData[year]
			rows = append(rows, []string{
				strconv.Itoa(year),
				formatNumber(stats.TotalDifferentials),
				formatNumber(stats.TotalLinesChanged),
			})
		}

		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Background(lipgloss.Color(_background))).
			Headers(headers...).
			Rows(rows...).
			Width(60)

		fmt.Fprint(w, t.Render(), "\n")
	}
}

func New() entity.Reporter {
	return &reporter{}
}

var _ entity.Reporter = &reporter{}