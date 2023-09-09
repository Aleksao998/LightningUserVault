package server

import (
	"fmt"

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

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.logLevelRaw,
		logLevelFlag,
		"WARN",
		"the log level for console output",
	)

	cmd.Flags().StringVar(
		&params.serverAddressRaw,
		serverAddressFlag,
		fmt.Sprintf("%s:%s", helper.DefaultServerEndpoint, helper.DefaultServerPort),
		"server endpoint",
	)

	cmd.Flags().BoolVar(
		&params.enableCache,
		enabledCacheFlag,
		true,
		"flag which represents if cache mechanism is enabled",
	)

	cmd.Flags().StringVar(
		&params.cacheTypeRaw,
		cacheTypeFlag,
		"MEMCACHE",
		"the type of cache, supported [MEMCACHE]",
	)

	cmd.Flags().StringVar(
		&params.memcacheAddressRaw,
		memcacheAddressFlag,
		fmt.Sprintf("%s:%s", helper.DefaultServerEndpoint, helper.DefaultMemcachePort),
		"memcache endpoint",
	)

	cmd.Flags().StringVar(
		&params.storageTypeRaw,
		storageTypeFlag,
		"PEBBLE",
		"the type of storage, supported [PEBBLE, POSTRESQL]",
	)

	cmd.Flags().StringVar(
		&params.dbHostRaw,
		dbHostRawFlag,
		fmt.Sprintf("%s:%s", helper.DefaultServerEndpoint, helper.DefaultDatabasePort),
		"database host endpoint",
	)

	cmd.Flags().StringVar(
		&params.dbUser,
		dbUserFlag,
		"postgres",
		"database user",
	)

	cmd.Flags().StringVar(
		&params.dbPass,
		dbPassFlag,
		"postgres",
		"database password",
	)

	cmd.Flags().StringVar(
		&params.dbName,
		dbNameFlag,
		"postgres",
		"database name",
	)
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter()
	if err := runServerLoop(outputter); err != nil {
		outputter.SetError(err)
		outputter.WriteOutput()

		return
	}
}

func runPreRun(cmd *cobra.Command, _ []string) error {
	return params.initRawParams()
}

func runServerLoop(
	outputter command.OutputFormatter,
) error {
	serverInstance, err := server.NewServer(params.generateConfig())
	if err != nil {
		return err
	}

	return helper.HandleSignals(serverInstance.Close, outputter)
}
