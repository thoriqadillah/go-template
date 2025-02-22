package command

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/stephenafamo/bob"
)

type App struct {
	BobDB *bob.DB
}

type CommandFunc func(ctx context.Context, app *App) *cobra.Command

var cmds = make([]CommandFunc, 0)

func Register(cmd ...CommandFunc) {
	cmds = append(cmds, cmd...)
}

var root = &cobra.Command{
	Use:   "cmd",
	Short: "Command line for the application",
}

func Execute(ctx context.Context, app *App) error {
	for _, cmd := range cmds {
		root.AddCommand(cmd(ctx, app))
	}

	return root.Execute()
}
