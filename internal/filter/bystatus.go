package filter

import "github.com/moisesvega/diffy/internal/entity"

// ByStatus returns a slices.DeleteFunc that filters differentials by status.
func ByStatus(status entity.Status) func(differential *entity.Differential) bool {
	return func(differential *entity.Differential) bool {
		return differential.Status != status
	}
}
