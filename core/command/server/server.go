package server

import (
	"fmt"
	"os"

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
		getEnvWithDefault("LOG_LEVEL", "WARN"),
		"the log level for console output",
	)

	cmd.Flags().StringVar(
		&params.serverAddressRaw,
		serverAddressFlag,
		getEnvWithDefault("SERVER_ADDRESS", fmt.Sprintf("%s:%s", helper.DefaultServerEndpoint, helper.DefaultServerPort)),
		"server endpoint",
	)

	cmd.Flags().StringVar(
		&params.enableCache,
		enabledCacheFlag,
		getEnvWithDefault("ENABLE_CACHE", "true"),
		"flag which represents if cache mechanism is enabled",
	)

	cmd.Flags().StringVar(
		&params.cacheTypeRaw,
		cacheTypeFlag,
		getEnvWithDefault("CACHE_TYPE", "MEMCACHE"),
		"the type of cache, supported [MEMCACHE]",
	)

	cmd.Flags().StringVar(
		&params.memcacheAddressRaw,
		memcacheAddressFlag,
		getEnvWithDefault("MEMCACHE_ADDRESS", fmt.Sprintf("%s:%s", helper.DefaultServerEndpoint, helper.DefaultMemcachePort)),
		"memcache endpoint",
	)

	cmd.Flags().StringVar(
		&params.storageTypeRaw,
		storageTypeFlag,
		getEnvWithDefault("STORAGE_TYPE", "PEBBLE"),
		"the type of storage, supported [PEBBLE, POSTRESQL]",
	)

	cmd.Flags().StringVar(
		&params.dbHostRaw,
		dbHostRawFlag,
		getEnvWithDefault("DB_HOST", fmt.Sprintf("%s:%s", helper.DefaultServerEndpoint, helper.DefaultDatabasePort)),
		"database host endpoint",
	)

	cmd.Flags().StringVar(
		&params.dbUser,
		dbUserFlag,
		getEnvWithDefault("DB_USER", "postgres"),
		"database user",
	)

	cmd.Flags().StringVar(
		&params.dbPass,
		dbPassFlag,
		getEnvWithDefault("DB_PASS", "postgres"),
		"database password",
	)

	cmd.Flags().StringVar(
		&params.dbName,
		dbNameFlag,
		getEnvWithDefault("DB_NAME", "postgres"),
		"database name",
	)
}

func runCommand(cmd *cobra.Command, _ []string) {
	if err := runServerLoop(); err != nil {
		fmt.Println("ERROR: ", err)

		return
	}
}

func runPreRun(cmd *cobra.Command, _ []string) error {
	return params.initRawParams()
}

func runServerLoop() error {
	serverInstance, err := server.NewServer(params.generateConfig())
	if err != nil {
		return err
	}

	return helper.HandleSignals(serverInstance.Close)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		fmt.Println(fmt.Sprintf("Exists %s", value))

		return value
	}

	return defaultValue
}
