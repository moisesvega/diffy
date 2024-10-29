package filter

import "github.com/moisesvega/diffy/internal/entity"

// MinLineCount returns a slices.DeleteFunc that filters differentials by line count.
func MinLineCount(count int) func(differential *entity.Differential) bool {
	return func(differential *entity.Differential) bool {
		return differential.LineCount < count
	}
}
