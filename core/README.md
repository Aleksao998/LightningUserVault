<h1 align="center"> Core </h1>

<p align="center"> The heart of the application, containing the main logic and components. </p>

## Components:

📡 **Server:** Currently HTTP-based, with plans to support WebSockets and gRPC

🚦 **Router:** Directs incoming requests to appropriate endpoints.

🛠️ **Handlers:** Process requests, interacting with cache and storage

⚡ **Caching Mechanism:** Ensures fast data retrieval

🗃️ **Storages:** Handles data persistence, supporting both Pebble and PostgreSQL

## Additional Information:

🖥️ **CLI:** Uses Cobra commands

🧪 **Testing Framework:** End-to-end testing framework for the HTTP server, offering utilities for server management during tests

## Mocking:
Each module in the core has a mock version for testing purposes. This includes database clients and cache clients