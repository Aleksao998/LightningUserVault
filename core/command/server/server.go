package server

import (
	"github.com/Aleksao998/LightingUserVault/core/command"
	"github.com/Aleksao998/LightingUserVault/core/command/helper"
	"github.com/Aleksao998/LightingUserVault/core/server"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:     "server",
		Short:   "The default command that starts LightingUserVault",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(serverCmd)

	return serverCmd
}

func setFlags(cmd *cobra.Command) {}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter()
	if err := runServerLoop(params.generateConfig(), outputter); err != nil {
		outputter.SetError(err)
		outputter.WriteOutput()

		return
	}
}

func runPreRun(cmd *cobra.Command, _ []string) error {
	return params.initRawParams()
}

func runServerLoop(
	config *server.Config,
	outputter command.OutputFormatter,
) error {
	serverInstance, err := server.NewServer(config)
	if err != nil {
		return err
	}

	return helper.HandleSignals(serverInstance.Close, outputter)
}
