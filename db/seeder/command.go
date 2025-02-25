package seeder

import (
	"app/cmd/command"
	"context"

	"github.com/spf13/cobra"
)

func CreateCommand(ctx context.Context, app *command.App) *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Seed the database",
		Run: func(cmd *cobra.Command, args []string) {
			if err := Seed(ctx, app.Db); err != nil {
				panic(err)
			}
		},
	}
}
