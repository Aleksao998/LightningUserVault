package server

import (
	"net"

	"github.com/Aleksao998/LightningUserVault/core/command/server/types"
	"go.uber.org/zap/zapcore"
)

// Config is used to parametrize the LightningUserVault client
type Config struct {
	// LogLevel is a log type
	LogLevel zapcore.Level

	// ServerAddress is an address of http server
	ServerAddress *net.TCPAddr

	// EnableCache is a flag which represents if cache mechanism is enabled
	EnableCache bool

	// CacheType is a cache type [MEMCACHE]
	CacheType types.CacheType

	// MemcacheAddress is an address of memcache server
	MemcacheAddress *net.TCPAddr

	// StorageType is a cache type [PEBBLE, POSTRESQL]
	StorageType types.StorageType

	// DBHost is an address of database host
	DBHost *net.TCPAddr

	// DBUser is a user name for database
	DBUser string

	// DBPass is a password for database user
	DBPass string

	// DBName is a name of database
	DBName string
}
