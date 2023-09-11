package server

import (
	"net"
	"testing"

	"github.com/Aleksao998/LightningUserVault/core/command/server/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestInitRawParams(t *testing.T) {
	t.Parallel()

	sp := &serverParams{
		logLevelRaw:        "DEBUG",
		serverAddressRaw:   "localhost:8080",
		cacheTypeRaw:       "MEMCACHE",
		memcacheAddressRaw: "localhost:11211",
		storageTypeRaw:     "PEBBLE",
		dbHostRaw:          "localhost:5432",
	}

	err := sp.initRawParams()
	assert.NoError(t, err)
	assert.Equal(t, zapcore.DebugLevel, sp.logLevel)
	assert.NotNil(t, sp.serverAddress)
	assert.NotNil(t, sp.memcacheAddress)
	assert.NotNil(t, sp.dbHost)
}

func TestGenerateConfig(t *testing.T) {
	sp := &serverParams{
		logLevel:        zapcore.DebugLevel,
		serverAddress:   &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080},
		enableCache:     "true",
		cacheType:       types.MEMCACHE,
		memcacheAddress: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 11211},
		storageType:     types.PEBBLE,
		dbHost:          &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5432},
		dbUser:          "user",
		dbPass:          "pass",
		dbName:          "testdb",
	}

	config := sp.generateConfig()
	assert.Equal(t, zapcore.DebugLevel, config.LogLevel)
	assert.Equal(t, sp.serverAddress, config.ServerAddress)
	assert.Equal(t, sp.cacheType, config.CacheType)
	assert.Equal(t, sp.memcacheAddress, config.MemcacheAddress)
	assert.Equal(t, sp.storageType, config.StorageType)
	assert.Equal(t, sp.dbHost, config.DBHost)
	assert.Equal(t, sp.dbUser, config.DBUser)
	assert.Equal(t, sp.dbPass, config.DBPass)
	assert.Equal(t, sp.dbName, config.DBName)
}
