package cmd

import (
	"time"

	"github.com/moisesvega/diffy/internal/model"
	"github.com/moisesvega/diffy/internal/reporter/heatmap"
	"github.com/spf13/cobra"
)

func Main() *cobra.Command {
	o := &opts{}
	cmd := &cobra.Command{
		Use:           "diffy",
		Short:         "CLI designed to deliver comprehensive statistics and insights from code reviews and differential analysis",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	setFlags(cmd.Flags(), o)

	hm := heatmap.New()
	hm.Report([]model.User{
		{
			Username: "moisesvega",
			Differentials: []*model.Differential{
				{
					ID:             "1",
					Title:          "New",
					LineCount:      "10",
					Status:         "Nice",
					URI:            "ww.example.com",
					CreatedAt:      time.Now(),
					ModifiedAt:     time.Now(),
					RepositoryPHID: "",
				},
			},
		},
	})

	return cmd
}
