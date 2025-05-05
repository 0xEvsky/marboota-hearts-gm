# Seat

## Overview

The `Seat` class represents a seating position at the game table. It manages the relationship between a physical position in the game and a player entity, handling the seating and unseating processes.

**File:** `scripts/seat_1.gd`

## Key Responsibilities

1. Managing the association between seat position and player
2. Handling seat interaction (clicking to sit)
3. Communicating with the server about seating changes
4. Controlling seat availability

## Properties

| Property | Type | Description |
|----------|------|-------------|
| seat_num | Integer | Unique identifier for this seat position |
| sitter | Player | Reference to the player occupying this seat (or null) |
| is_taken | Boolean | Whether the seat is currently occupied |

## Node References

| Node | Type | Description |
|------|------|-------------|
| seat_ready_button | Button | Reference to the shared ready button |

## Main Functions

### `_disable_button()` / `_enable_button()`
Controls the interactivity of the seat button based on occupancy.

### `seat_player(id)`
Assigns a player to this seat:
1. Retrieves the player instance from the PlayerManager
2. If player is already seated elsewhere, unseats them first
3. Updates player state to PLAYER_WAITING
4. Moves player to seat position
5. Updates references and disables the button

### `unseat_player()`
Removes a player from this seat:
1. Calls the player's unseat method
2. Clears the sitter reference
3. Re-enables the button

### `_on_button_button_up()`
Button click handler that:
1. Attempts to seat the local player
2. Shows the ready button
3. Sends a sit request to the server
4. Handles potential errors by reverting changes

## Integration with Other Components

- **Table**: Parent node that contains all seats
- **PlayerManager**: Used to access player instances
- **EventManager**: Used to send sit requests to the server
- **Player**: Linked to seat when occupied
- **ReadyButton**: Controlled by seat when local player sits

## Seating Process

1. Player clicks on an available seat
2. Local player is visually positioned at the seat
3. Seat request is sent to the server
4. If successful, nothing changes
5. If unsuccessful, player is returned to pinned players list

## Example Usage

```gdscript
# Accessing a seat
var seat1 = get_node("Table/Seat1")

# Checking if seat is occupied
if seat1.is_taken:
    print("Seat is taken by: ", seat1.sitter.username)
    
# Programmatically seating a player
seat1.seat_player("player_id")

# Unseating a player
seat1.unseat_player()
```

## Seat Layout

In the current implementation, the game table has 4 seats arranged in a cross pattern:
- Seat1: Bottom (0, 313)
- Seat2: Right (313, 0)
- Seat3: Top (0, -313)
- Seat4: Left (-313, 0)

This layout creates a 4-player game table with positions for players around the central table area.

## Future Improvements

1. Add seat animations for transitions
2. Implement different seat types/classes
3. Add support for team-based seating
4. Include player name tags above seats
