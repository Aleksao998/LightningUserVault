package framework

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Aleksao998/LightingUserVault/core/command/helper"
	serverCommand "github.com/Aleksao998/LightingUserVault/core/command/server"
	"github.com/Aleksao998/LightingUserVault/core/server"
	"go.uber.org/zap/zapcore"
)

const (
	initialPort = 12000
	localhost   = "localhost"
)

// TestServer represents a server used for testing purposes
type TestServer struct {
	t      *testing.T
	Config *server.Config
	Port   *ReservedPort
	cmd    *exec.Cmd
}

// NewTestServer initializes a new TestServer instance
func NewTestServer(t *testing.T) *TestServer {
	t.Helper()

	port := FindAvailablePort(initialPort, initialPort+10000)

	host := localhost + ":" + port.Port()

	serverAddress, err := helper.ResolveAddr(
		host,
		helper.LocalHostBinding,
	)
	if err != nil {
		t.Fatal(err)
	}

	config := server.Config{
		ServerAddress: serverAddress,
		EnableCache:   false,
		LogLevel:      zapcore.DebugLevel,
	}

	return &TestServer{
		t:      t,
		Config: &config,
		Port:   port,
	}
}

// Stop terminates the test server process
func (t *TestServer) Stop() {
	if t.cmd != nil {
		// Send a SIGINT signal to the process to trigger a graceful shutdown
		if err := t.cmd.Process.Signal(os.Interrupt); err != nil {
			t.t.Error(err)
		}
	}
}

// ReleaseReservedPorts releases any reserved ports associated with the test server
func (t *TestServer) ReleaseReservedPorts() {
	if err := t.Port.Close(); err != nil {
		t.t.Error(err)
	}

	t.Port = nil
}

// Start initiates the test server with the given context
func (t *TestServer) Start(ctx context.Context) error {
	serverCmd := serverCommand.GetCommand()

	args := []string{
		serverCmd.Use,
		"--server-address", t.Config.ServerAddress.String(),
		"--enable-cache", strconv.FormatBool(t.Config.EnableCache),
		"--log-level", t.Config.LogLevel.String(),
	}
	fmt.Println(args)
	t.ReleaseReservedPorts()

	binaryName := os.Getenv("BINARY_PATH")
	if binaryName == "" {
		binaryName = "build/lighting_user_vault" // default value
	}

	t.cmd = exec.Command(binaryName, args...)

	stdout := io.Writer(os.Stdout)
	t.cmd.Stdout = stdout
	t.cmd.Stderr = stdout

	if err := t.cmd.Start(); err != nil {
		return err
	}

	// In the future we should not be dependent on timeout
	// to know that server started. We should have pooling mechanism which validates
	// that server is up
	time.Sleep(5 * time.Second)

	return nil
}

// NewTestServerAndStart initializes and starts a new TestServer instance in a separate goroutine
func NewTestServerAndStart(t *testing.T) *TestServer {
	t.Helper()

	srv := NewTestServer(t)

	t.Cleanup(func() {
		srv.Stop()
	})

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		if err := srv.Start(ctx); err != nil {
			t.Fatal("server failed to start", err)
		}
	}()

	wg.Wait()

	return srv
}

func CleanupStorage() {
	os.RemoveAll("pebble-storage")
}
