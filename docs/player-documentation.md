# Player

## Overview

The `Player` class represents an individual player entity in the game. Each player has a visual representation, state management, and seat association.

**File:** `scripts/player.gd`

## Key Responsibilities

1. Representing a player visually with an icon
2. Tracking player state (idle, waiting, ready, etc.)
3. Managing seat association
4. Handling the unseat process

## Properties

| Property | Type | Description |
|----------|------|-------------|
| username | String | Player's display name |
| state | Integer | Current player state (from PlayerManager enum) |
| seat | Seat | Reference to the seat this player is occupying (or null) |

## Node References

| Node | Type | Description |
|------|------|-------------|
| manager | PlayerManager | Reference to the parent PlayerManager |
| icon | Sprite2D | The visual representation of the player |

## Main Functions

### `unseat()`
Handles the process of a player standing up from a seat:
1. Changes player state to PLAYER_IDLE
2. Clears seat reference
3. Moves player back to the player list (pinned players)

## Player States

The Player class uses the states defined in PlayerManager:

| State | Description |
|-------|-------------|
| PLAYER_UNAVAILABLE | Player is not available for interaction |
| PLAYER_IDLE | Player is connected but not engaged in gameplay |
| PLAYER_WAITING | Player is seated but not ready |
| PLAYER_READY | Player is seated and ready to play |
| PLAYER_TRUMPING | Player is in trump selection phase (game-specific) |
| PLAYER_PLAYING | Player is actively playing |

## Integration with Other Components

- **PlayerManager**: Creates and manages Player instances
- **Seat**: Associates with a Player when occupied
- **Globals**: Local player is accessible via Globals.my_player
- **ReadyButton**: Modifies player state when toggled

## Visual Representation

The Player scene includes:
1. An IconBorder Sprite2D (slightly larger than the icon)
2. An Icon Sprite2D (the player's avatar)

The icon texture is loaded dynamically from the URL provided during player creation.

## Local vs. Remote Players

While the Player class handles both local and remote players the same way, the local player (the one controlled by the user) is specially tracked:

1. It's accessible globally via `Globals.my_player`
2. It's the only player that can send requests directly (sit, ready, etc.)

## Example Usage

```gdscript
# Accessing the local player
var my_player = Globals.my_player

# Changing player state
my_player.state = Globals.player_manager.PLAYER_READY

# Checking if player is seated
if my_player.seat != null:
    print("Player is seated at position: ", my_player.global_position)
    
# Making a player stand up
player.unseat()
```

## Future Improvements

1. Add player customization options
2. Implement player animation
3. Add player emotes or reactions
4. Add player statistics tracking
