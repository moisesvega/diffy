package heatmap

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/moisesvega/diffy/internal/model"
)

const (
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

func (r *reporter) Report(users []model.User) error {

	heatmap := make(map[time.Time]int)
	for _, user := range users {
		for _, differential := range user.Differentials {
			heatmap[differential.ModifiedAt]++
		}
	}

	// TODO(moisesvega): Pass an option to set the date range
	var since *time.Time
	what := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local)
	// what := time.Now().AddDate(-1, 0, 0)
	since = &what

	today := time.Now().AddDate(0, 0, -1)
	yearAgo := today.AddDate(-1, 0, 0)
	if since != nil {
		yearAgo = *since
	}

	for yearAgo.Weekday() != time.Saturday {
		yearAgo = yearAgo.AddDate(0, 0, +1)
	}

	// 7 days of the week + 1 row for the month
	rows := make([][]string, 7)
	// First column should be the days of the week
	for weekday := range 7 {
		rows[weekday] = append(rows[weekday], time.Weekday(weekday).String())
	}

	// We need to track the current month to add it to the heatmap
	currentMonth := yearAgo.Month()
	// To make it easier to add the month to the heatmap, we'll add it to the last row
	headers := make([]string, 0)
	headers = append(headers, "", currentMonth.String()[0:3])
	total := 0
	for !yearAgo.After(today) {
		yearAgo = yearAgo.AddDate(0, 0, 1)
		count := rand.IntN(10)
		if v, ok := heatmap[yearAgo]; ok {
			count = v
		}
		diffCount := strconv.Itoa(count)
		if count == 0 {
			diffCount = ""
		}
		rows[yearAgo.Weekday()] = append(rows[yearAgo.Weekday()], diffCount)
		if yearAgo.Month() != currentMonth {
			currentMonth = yearAgo.Month()
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
			fontColor := _blackFontColor
			switch {
			case v >= 1 && v < 3:
				color = _low
				fontColor = _grayFontColor
			case v >= 3 && v < 5:
				color = _mid
				fontColor = _grayFontColor
			case v >= 5 && v < 7:
				color = _high
			case v >= 7:
				color = _max
			}

			return lipgloss.NewStyle().
				Background(lipgloss.Color(color)).
				Width(3).
				Foreground(lipgloss.Color(fontColor)).
				Align(lipgloss.Center).BorderBackground(lipgloss.Color("#57606a"))
		}).
		Rows(rows...).Width(0).Headers(headers...)
	fmt.Println(t)
	fmt.Println("Total Differentials:" + strconv.Itoa(total))
	return nil
}

func New() model.Reporter {
	return &reporter{}
}

var _ model.Reporter = &reporter{}
