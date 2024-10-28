package filter

import "github.com/moisesvega/diffy/internal/model"

// ByLineCount returns a slices.DeleteFunc that filters differentials by status.
func ByLineCount(count int) func(differential *model.Differential) bool {
	return func(differential *model.Differential) bool {
		return differential.LineCount < count
	}
}
