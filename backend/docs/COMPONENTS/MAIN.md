# COMPONENTS/MAIN.md

## Main Server (`main.go`)

### Purpose:
- Initializes the server application.
- Sets up HTTP routes and WebSocket endpoints.

### Key Responsibilities:
- Bootstrapping application configuration.
- Starting the WebSocket service.

### Interactions:
- Calls `ws.HandleConnections()` to manage WebSocket logic.
- Acts as a centralized entry point.
