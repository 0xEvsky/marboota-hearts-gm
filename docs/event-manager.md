# Event Manager

## Overview

The `EventManager` is a singleton autoload script that handles game event processing and message dispatching. It serves as a broker between the NetworkManager and game components, translating network messages into game events.

**File:** `scripts/event_manager.gd`

## Key Responsibilities

1. Processing incoming network messages
2. Emitting signals for game events
3. Managing request/response communication with the server
4. Handling success and error callbacks for requests
5. Creating properly formatted game action requests

## Signals

| Signal | Description |
|--------|-------------|
| JOIN_received | Emitted when a new player joins the game |
| LEAVE_received | Emitted when a player leaves the game |
| SIT_received | Emitted when a player sits at a seat |
| UNSIT_received | Emitted when a player stands up from a seat |
| READY_received | Emitted when a player marks themselves as ready |
| UNREADY_received | Emitted when a player marks themselves as not ready |

## Main Functions

### `send_request(msg, on_success, on_error)`
Sends a request to the server and registers callback functions for success and error responses.

Parameters:
- `msg`: Dictionary containing the request message
- `on_success`: Callable function to execute on success
- `on_error`: Callable function to execute on error

### `_handle_message(msg)`
Processes incoming messages from NetworkManager, determining if they are responses to requests or new events.

### `_dispatch(action, msg)`
Emits the appropriate signal based on the action type in the message.

### Request Helper Functions

The EventManager provides convenience functions for creating properly formatted requests:

| Function | Description |
|----------|-------------|
| `sit_request(seat)` | Creates request to sit at specified seat |
| `unsit_request()` | Creates request to stand up from current seat |
| `ready_request()` | Creates request to mark player as ready |
| `unready_request()` | Creates request to mark player as not ready |

## Request/Response Flow

1. A game component calls `send_request()` with a message and callbacks
2. EventManager adds a unique request ID to the message
3. Message is sent to server via NetworkManager
4. When a response with matching ID returns, appropriate callback is executed
5. For non-request messages, signals are emitted to notify subscribers

## Integration with Other Components

- **NetworkManager**: Source of raw messages for processing
- **PlayerManager**: Subscribes to JOIN, LEAVE, and SIT signals
- **Seat**: Uses request functions to communicate sit/unsit actions
- **ReadyButton**: Uses request functions to communicate ready/unready actions

## Example Usage

```gdscript
# Send a sit request
EventManager.send_request(
    EventManager.sit_request(1),  # Sit at seat 1
    func():  # Success callback
        print("Successfully sat down"),
    func(error):  # Error callback
        print("Failed to sit down: ", error)
)

# Subscribe to an event
EventManager.JOIN_received.connect(func(user_id, username, icon_url):
    print("Player joined: ", username)
)
```

## Future Improvements

1. Implement message validation
2. Add timeout handling for requests
3. Create a message queue for rate limiting
4. Add better error handling and reporting
