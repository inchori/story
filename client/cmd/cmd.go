// Package cmd provides the cli for running the iliad consensus client.
package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/storyprotocol/iliad/client/app"
	iliadcfg "github.com/storyprotocol/iliad/client/config"
	"github.com/storyprotocol/iliad/lib/buildinfo"
	libcmd "github.com/storyprotocol/iliad/lib/cmd"
	"github.com/storyprotocol/iliad/lib/log"
)

// New returns a new root cobra command that handles our command line tool.
func New() *cobra.Command {
	return libcmd.NewRootCmd(
		"iliad",
		"Iliad is a consensus client implementation for the Story L1 blockchain",
		newRunCmd("run", app.Run),
		newInitCmd(),
		buildinfo.NewVersionCmd(),
		newValidatorCmds(),
	)
}

// newRunCmd returns a new cobra command that runs the iliad consensus client.
func newRunCmd(name string, runFunc func(context.Context, app.Config) error) *cobra.Command {
	iliadCfg := iliadcfg.DefaultConfig()
	logCfg := log.DefaultConfig()

	cmd := &cobra.Command{
		Use:   name,
		Short: "Runs the iliad consensus client",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := log.Init(cmd.Context(), logCfg)
			if err != nil {
				return err
			}
			if err := libcmd.LogFlags(ctx, cmd.Flags()); err != nil {
				return err
			}

			cometCfg, err := parseCometConfig(ctx, iliadCfg.HomeDir)
			if err != nil {
				return err
			}

			return runFunc(ctx, app.Config{
				Config: iliadCfg,
				Comet:  cometCfg,
			})
		},
	}

	bindRunFlags(cmd, &iliadCfg)
	log.BindFlags(cmd.Flags(), &logCfg)

	return cmd
}
