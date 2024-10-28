package heatmap

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/moisesvega/diffy/internal/model"
)

const (
	// TODO(moisesvega): Make it configurable
	_zero           = "#11141A"
	_low            = "#0e4429"
	_mid            = "#006d32"
	_high           = "#26a641"
	_max            = "#39d353"
	_background     = "#212830"
	_whiteFontColor = "#FAFAFA"
	_grayFontColor  = "#D1D1D1"
	_blackFontColor = "#000000"
)

var style = lipgloss.NewStyle().
	Foreground(lipgloss.Color(_whiteFontColor)).
	Background(lipgloss.Color(_background)).Align(lipgloss.Center)

type reporter struct {
}

const _timeLayout = "2006-01-02"

func (r *reporter) Report(users []*model.User, option ...model.ReporterOption) error {
	opts := &model.ReporterOptions{}
	// Apply the options
	for _, o := range option {
		o(opts)
	}
	for _, user := range users {
		reportUser(user, opts)
	}
	return nil
}

func reportUser(user *model.User, opts *model.ReporterOptions) {
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}
	heatmap := make(map[string]int)
	for _, differential := range user.Differentials {
		heatmap[differential.ModifiedAt.Format(_timeLayout)]++
	}

	today := time.Now().AddDate(0, 0, -1)
	// Default since is the beginning of the year
	since := time.Date(
		time.Now().Year(),
		time.January, 1, 0, 0, 0, 0,
		time.Local).AddDate(-1, 0, 0)

	// If the user provided a "since" date, we'll use that instead
	if opts.Since != nil {
		since = *opts.Since
	}

	// We need to find the first Sunday of the year
	for since.Weekday() != time.Sunday {
		since = since.AddDate(0, 0, +1)
	}

	// 7 days of the week + 1 row for the month
	rows := make([][]string, 7)
	// First column should be the days of the week
	for weekday := range 7 {
		rows[weekday] = append(rows[weekday], time.Weekday(weekday).String())
	}

	// We need to track the current month to add it to the heatmap
	currentMonth := since.Month()
	// To make it easier to add the month to the heatmap, we'll add it to the last row
	headers := make([]string, 0)
	headers = append(headers, "", currentMonth.String()[0:3])
	total := 0
	for !since.After(today) {
		since = since.AddDate(0, 0, 1)
		count := 0
		if v, ok := heatmap[since.Format(_timeLayout)]; ok {
			count = v
		}
		diffCount := strconv.Itoa(count)
		rows[since.Weekday()] = append(rows[since.Weekday()], diffCount)
		if since.Month() != currentMonth {
			currentMonth = since.Month()
			for len(headers)-len(rows[0]) <= 0 {
				headers = append(headers, "")
			}
			headers = append(headers, currentMonth.String()[0:3])
		}
		total += count
	}

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderStyle(lipgloss.NewStyle().Background(lipgloss.Color(_background)).Width(0)).
		StyleFunc(func(row, col int) lipgloss.Style {
			rowFromData := row - 1
			if row > len(rows) || rowFromData < 0 {
				return style
			}
			if col >= len(rows[rowFromData]) {
				return style
			}

			value := rows[rowFromData][col]
			v, err := strconv.Atoi(value)
			if err != nil {
				return style
			}

			color := _zero
			fontColor := _whiteFontColor
			switch {
			case v >= 1 && v < 3:
				color = _low
				fontColor = _grayFontColor
			case v >= 3 && v < 5:
				color = _mid
				fontColor = _grayFontColor
			case v >= 5 && v < 7:
				color = _high
				fontColor = _blackFontColor
			case v >= 7:
				color = _max
				fontColor = _blackFontColor
			}

			return lipgloss.NewStyle().
				Background(lipgloss.Color(color)).
				Width(3).
				Foreground(lipgloss.Color(fontColor)).
				Align(lipgloss.Center).BorderBackground(lipgloss.Color("#57606a"))
		}).
		Rows(rows...).Width(0).Headers(headers...)
	_, _ = fmt.Fprint(w, t.Render())
}

func New() model.Reporter {
	return &reporter{}
}

var _ model.Reporter = &reporter{}
