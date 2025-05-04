# Ready System Documentation

## Overview

The Ready System allows players to indicate they are prepared to start a game session. This functionality is essential for synchronizing player readiness in a multiplayer environment and ensuring all participants are prepared before gameplay begins.

## Core Components

### Ready Button (`ready_button.gd`)

The primary interface element is a toggle button that:
- Appears when a player sits at a table
- Can be toggled between ready/not ready states
- Uses visual styling to indicate the current state
- Communicates with the server to sync ready status

## Button States

The Ready Button has two visual states:
1. **Normal state**: Indicates player is seated but not ready
2. **Pressed state**: Indicates player is ready (styled with a green background)

## Ready Process

### Local Player Ready

When a local player toggles the Ready button to ON:

1. Player's state is optimistically updated to `PLAYER_READY`
2. A network request is sent via `EventManager.ready_request()`
3. On success, the player remains in ready state
4. On failure, the player's state is reverted to `PLAYER_WAITING`

### Local Player Unready

When a local player toggles the Ready button to OFF:

1. Player's state is optimistically updated to `PLAYER_WAITING`
2. A network request is sent via `EventManager.unready_request()`
3. On success, the player remains in unready state
4. On failure, the player's state is reverted to `PLAYER_READY`

### Remote Player Ready/Unready

When a remote player changes ready status:
1. The server broadcasts a READY or UNREADY event
2. `EventManager` emits the corresponding signal (`READY_received` or `UNREADY_received`)
3. Other components can listen for these signals to update the UI accordingly

## Button Visibility

The Ready button is initially hidden and appears only when:
- The local player sits at a seat
- The button is hidden again when the player stands up

## Integration with Game Flow

The Ready System integrates with the overall game flow:

1. Players join the game and sit at the table
2. Players toggle their ready status using the Ready button
3. When all players are ready, the game can progress to the next phase
   (Note: The actual game start logic is not fully implemented in the provided code)

## Network Communication

Ready state changes follow this pattern:

```
[Local UI Change] → [Send Request] → [Wait for Response] → [Keep or Revert Change]
```

The request/response flow:

```
// Client -> Server (Ready Request)
{
  "ACTION": "READY",
  "REQUESTID": "request-12345-789"
}

// Server -> Client (Success)
{
  "ACTION": "OK",
  "REQUESTID": "request-12345-789"
}

// Server -> All Clients (Broadcast)
{
  "ACTION": "READY",
  "USERID": "12345"
}
```

## Visual Feedback

The Ready button uses:
- Toggle mode to maintain pressed/unpressed state
- Custom styling for the pressed state (green background)
- Text label "Ready" to indicate its purpose

## Code Implementation

The core logic in the `_on_toggled` method handles both ready and unready actions:

```gdscript
func _on_toggled(toggled_on: bool) -> void:
    if toggled_on:
        // Ready functionality
        Globals.my_player.state = Globals.player_manager.PLAYER_READY
        EventManager.send_request(EventManager.ready_request(),
            func(): pass,  // Success handler
            func():        // Error handler
                Globals.my_player.state = Globals.player_manager.PLAYER_WAITING
        )
    else:
        // Unready functionality
        if Globals.my_player.state == Globals.player_manager.PLAYER_READY:
            Globals.my_player.state = Globals.player_manager.PLAYER_WAITING
            EventManager.send_request(EventManager.unready_request(),
                func(): pass,  // Success handler
                func():        // Error handler
                    Globals.my_player.state = Globals.player_manager.PLAYER_READY
            )
```

## Limitations and Considerations

- The current implementation focuses on the UI and network communication
- There is no game-start logic implemented when all players are ready
- The system does not keep track of which players are ready (this would likely be handled by the server)
- There's no visual indication of other players' ready status