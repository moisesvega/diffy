package filter

import "github.com/moisesvega/diffy/internal/model"

// MinLineCount returns a slices.DeleteFunc that filters differentials by line count.
func MinLineCount(count int) func(differential *model.Differential) bool {
	return func(differential *model.Differential) bool {
		return differential.LineCount < count
	}
}
