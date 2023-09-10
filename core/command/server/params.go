package server

import (
	"log"
	"net"
	"strconv"

	"github.com/Aleksao998/LightingUserVault/core/command/helper"
	"github.com/Aleksao998/LightingUserVault/core/command/server/types"
	"github.com/Aleksao998/LightingUserVault/core/server"
	"go.uber.org/zap/zapcore"
)

var (
	params = &serverParams{}
)

const (
	logLevelFlag        = "log-level"
	serverAddressFlag   = "server-address"
	enabledCacheFlag    = "enable-cache"
	cacheTypeFlag       = "cache-type"
	memcacheAddressFlag = "memcache-address"
	storageTypeFlag     = "storage-type"
	dbHostRawFlag       = "database-host"
	dbUserFlag          = "database-user"
	dbPassFlag          = "database-pass"
	dbNameFlag          = "database-name"
)

type serverParams struct {
	// logLevel is a log type [ERROR, INFO, DEBUG, WARN]
	logLevel zapcore.Level

	// logLevelRaw is a raw log type
	logLevelRaw string

	// serverAddress is an address of http server
	serverAddress *net.TCPAddr

	// serverAddressRaw is a raw address of http server
	serverAddressRaw string

	// enableCache is a flag which represents if cache mechanism is enabled
	enableCache string

	// cacheType is a cache type [MEMCACHE]
	cacheType types.CacheType

	// cacheTypeRaw is a raw cache type
	cacheTypeRaw string

	// memcacheAddress is an address of memcache server
	memcacheAddress *net.TCPAddr

	// memcacheAddressRaw is a raw address of memcache server
	memcacheAddressRaw string

	// storageType is a cache type [PEBBLE, POSTGRESQL]
	storageType types.StorageType

	// storageTypeRaw is a raw storage type
	storageTypeRaw string

	// dbHost is an address of database host
	dbHost *net.TCPAddr

	// dbHostRaw is a raw address of database host
	dbHostRaw string

	// dbUser is a user name for database
	dbUser string

	// dbUser is a password for database user
	dbPass string

	// dbName is a name of database
	dbName string
}

func (p *serverParams) initRawParams() error {
	var err error

	// Parse log level
	p.logLevel, err = types.ConvertStringToLogLevel(p.logLevelRaw)
	if err != nil {
		return err
	}

	// Parse server address
	if p.serverAddress, err = helper.ResolveAddr(
		p.serverAddressRaw,
		helper.LocalHostBinding,
	); err != nil {
		return err
	}

	// Parse cache type
	p.cacheType, err = types.ConvertStringToCacheType(p.cacheTypeRaw)
	if err != nil {
		return err
	}

	// Parse memcache address
	if p.memcacheAddress, err = helper.ResolveAddr(
		p.memcacheAddressRaw,
		helper.LocalHostBinding,
	); err != nil {
		return err
	}

	// Parse storage type
	p.storageType, err = types.ConvertStringToStorageType(p.storageTypeRaw)
	if err != nil {
		return err
	}

	// Parse db host address
	if p.dbHost, err = helper.ResolveAddr(
		p.dbHostRaw,
		helper.LocalHostBinding,
	); err != nil {
		return err
	}

	return nil
}

func (p *serverParams) generateConfig() *server.Config {
	enableCache, err := strconv.ParseBool(p.enableCache)
	if err != nil {
		log.Fatal(err)
	}

	return &server.Config{
		LogLevel:        p.logLevel,
		ServerAddress:   p.serverAddress,
		EnableCache:     enableCache,
		CacheType:       p.cacheType,
		MemcacheAddress: p.memcacheAddress,
		StorageType:     p.storageType,
		DBHost:          p.dbHost,
		DBUser:          p.dbUser,
		DBPass:          p.dbPass,
		DBName:          p.dbName,
	}
}
