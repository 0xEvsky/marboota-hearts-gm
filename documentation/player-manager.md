# Player Manager

## Overview

The `PlayerManager` is responsible for creating, tracking, and managing all player instances in the game. It handles player joining, leaving, positioning, and maintains player state.

**File:** `scripts/player_manager.gd`

## Key Responsibilities

1. Creating and initializing player instances
2. Managing player positions and visibility
3. Handling player join/leave events
4. Processing player seat assignments
5. Maintaining a list of "pinned" players (players not sitting at the table)

## Properties

| Property | Type | Description |
|----------|------|-------------|
| pinned_players | Array[Player] | Players not seated at the table |

## Player States

The PlayerManager defines several player states as an enum:

| State | Description |
|-------|-------------|
| PLAYER_UNAVAILABLE | Player is not available for interaction |
| PLAYER_IDLE | Player is connected but not engaged in gameplay |
| PLAYER_WAITING | Player is seated but not ready |
| PLAYER_READY | Player is seated and ready to play |
| PLAYER_TRUMPING | Player is in trump selection phase (game-specific) |
| PLAYER_PLAYING | Player is actively playing |

## Main Functions

### `_ready()`
Initializes the manager, sets up signal connections, and registers itself with the Globals singleton.

### `_init_my_player()`
Initializes the local player after successful authentication.

### `_on_player_join(id, username, url)`
Creates a new player instance when a player joins the game.
- Downloads and applies the player's avatar image
- Adds the player to the pinned players list

### `_on_player_leave(id)`
Handles a player leaving the game.
- Removes player from their seat if they are seated
- Removes player from pinned players list
- Frees player instance

### `_on_player_sit(id, seat_num)`
Assigns a player to a specific seat.

### `move_player(id, pos)`
Updates a player's position in the scene.

### `pin_player(player)` / `unpin_player(player)`
Adds/removes a player from the pinned players list and updates positions.

## Integration with Other Components

- **NetworkManager**: PlayerManager listens for AUTH_accepted signal
- **EventManager**: Subscribes to JOIN_received, LEAVE_received, and SIT_received signals
- **Globals**: Registers itself with Globals for global access
- **Seat**: Seat components use PlayerManager to access player instances
- **Player**: Player instances are created and managed by PlayerManager

## Pinned Players System

The "pinned players" system manages players that are not seated at the table:

1. Players start as "pinned" when they join
2. When a player sits at a seat, they are "unpinned"
3. When a player stands up from a seat, they are "pinned" again
4. Pinned players are arranged vertically in the scene

## Avatar Loading

The PlayerManager handles fetching player avatars:

1. Creates an HTTPRequest for each player
2. Downloads the avatar image from the provided URL
3. Converts the downloaded data to a texture
4. Applies the texture to the player's icon

## Example Usage

```gdscript
# Getting a reference to the PlayerManager
var player_manager = Globals.player_manager

# Moving a player
player_manager.move_player("player_id", Vector2(100, 100))

# Accessing player states
if player.state == player_manager.PLAYER_READY:
    print("Player is ready!")
```

## Future Improvements

1. Add player teams functionality
2. Implement player score tracking
3. Add player persistence between game sessions
4. Create a more flexible positioning system for pinned players
