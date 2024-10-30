package contributions

import (
	"fmt"
	"slices"

	"github.com/moisesvega/diffy/internal/client/phabricator"
	"github.com/moisesvega/diffy/internal/entity"
	"github.com/moisesvega/diffy/internal/filter"
	"github.com/moisesvega/diffy/internal/reporter/heatmap"
	"github.com/spf13/cobra"
)

func New(client phabricator.Client) *cobra.Command {
	r := &runner{
		client:   client,
		reporter: heatmap.New(),
	}
	return &cobra.Command{
		Use:   "contributions",
		Short: "Get contributions from users",
		RunE:  r.runE,
	}
}

type runner struct {
	client   phabricator.Client
	reporter entity.Reporter
}

func (r *runner) runE(cmd *cobra.Command, args []string) error {
	if r.client == nil {
		return fmt.Errorf("phabricator client is not set")
	}
	u, err := r.client.GetUsers(args)
	if err != nil {
		return err
	}

	// filter out closed differentials and those with less than 10 lines
	for _, user := range u {
		diffs := slices.DeleteFunc(user.Differentials, filter.ByStatus(entity.Closed))
		user.Differentials = slices.DeleteFunc(diffs, filter.MinLineCount(10))
	}

	// report the data
	if err := r.reporter.Report(u); err != nil {
		return fmt.Errorf("failed to report: %w", err)
	}
	return nil
}
