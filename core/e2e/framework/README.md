## Overview
This e2e (end-to-end) framework is designed to facilitate the testing of our HTTP server. It provides utilities to start the server with default configurations and manage its lifecycle during tests.

## Features
1. **Server Initialization:** Easily start the HTTP server with default server configurations.
2. **Reserved Port Management:** The framework ensures that the server is started on a reserved port, avoiding conflicts with other running services.
3. **Test Server Management:** Provides functionalities to start, stop, and manage the test server's lifecycle.

## Limitations
1. **Database Support:** Currently, the e2e framework only supports the Pebble database. PostgreSQL is not supported in this version.
2. **Cache:** Caching is not supported in this version of the framework.

## Usage
🔴 Important: After each test, you need to clear the database to ensure a clean state for subsequent tests.

1. Starting the Server:
```go
testServer := freamwork.NewTestServerAndStart(t)
```

2. Stopping the Server:
```go
testServer.Stop()
```

3. Clearing the Database:
```go
freamwork.CleanupStorage()
```

## Future Enhancements
1. Support for additional databases like PostgreSQL.
2. Integration of caching mechanisms.