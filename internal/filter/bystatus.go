package filter

import "github.com/moisesvega/diffy/internal/model"

// ByStatus returns a slices.DeleteFunc that filters differentials by status.
func ByStatus(status model.Status) func(differential *model.Differential) bool {
	return func(differential *model.Differential) bool {
		return differential.Status != status
	}
}
