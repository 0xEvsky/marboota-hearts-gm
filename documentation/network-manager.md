# Network Manager

## Overview

The `NetworkManager` is a singleton autoload script that handles all communication with the backend server. It manages WebSocket connections, player authentication, and message exchange.

**File:** `scripts/network_manager.gd`

## Key Responsibilities

1. Establishing and maintaining WebSocket connection
2. Authenticating the player with the server
3. Reading and writing JSON messages
4. Processing incoming messages
5. Forwarding messages to the EventManager for further handling

## Properties

| Property | Type | Description |
|----------|------|-------------|
| instance_id | String | Unique instance identifier for the game |
| username | String | Player's display name |
| user_id | String | Unique player identifier |
| icon_url | String | URL to player's avatar image |
| _backend_url | String | WebSocket server address |
| _socket | WebSocketPeer | WebSocket connection object |
| authenticated | Boolean | Whether player is authenticated |

## Signals

| Signal | Description |
|--------|-------------|
| AUTH_accepted | Emitted when player has successfully authenticated with server |

## Main Functions

### `_ready()`
Initializes the WebSocket connection to the backend server.

### `_process(delta)`
Continuously polls the WebSocket connection, handles authentication, and processes incoming messages.

### `_handle_auth()`
Manages the authentication flow with the server, sending credentials and processing the response.

### `_read_loop()`
Reads incoming WebSocket messages and forwards them to the EventManager.

### `_write_json(msg)`
Serializes and sends JSON messages to the server.

### `_read_json()`
Deserializes incoming JSON messages from the server.

## Connection Flow

1. When the game starts, NetworkManager attempts connection to the WebSocket server
2. After connection, it sends an AUTH message with player credentials
3. If authentication is successful, it emits AUTH_accepted signal
4. The system then begins normal message processing

## Integration with Other Components

- **LoadingUI**: Hides when AUTH_accepted signal is emitted
- **EventManager**: Receives all incoming messages for further processing
- **PlayerManager**: Initializes when AUTH_accepted signal is emitted

## Example Usage

```gdscript
# Connect to authentication signal
NetworkManager.AUTH_accepted.connect(func():
    # Do something after authentication
    print("Player authenticated!")
)

# Access player information
print("Current player: ", NetworkManager.username)
```

## Future Improvements

1. Add reconnection logic for handling disconnections
2. Implement error handling for network failures
3. Add support for secure WebSocket connections (WSS)
4. Include timeout handling for connection attempts
