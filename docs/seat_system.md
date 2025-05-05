# Seat System Documentation

## Overview

The Seat System in Marboota Game manages how players interact with table positions. It handles the visual representation of seats, player seating/unseating, and coordinates with the network layer to ensure all clients stay synchronized.

## Core Components

### Seat Class (`seat_1.gd`)

Each seat is an instance of the `Seat` class with these key properties:

- `seat_num`: Unique identifier for the seat position (1-4)
- `sitter`: Reference to the Player object currently occupying the seat
- `is_taken`: Boolean indicating if the seat is occupied

### Table Scene (`table.tscn`)

The Table scene contains:
- 4 Seat instances positioned around a table sprite
- A ReadyButton that appears when a player sits down

## Seating Process

### Local Player Seating

When a local player clicks on a seat:

1. The seat's `_on_button_button_up()` method is called
2. The method calls `seat_player("Me")` to perform local seating
3. The Ready button is shown
4. A network request is sent via `EventManager.sit_request(seat_num)`
5. On success, the player remains seated
6. On failure:
   - The Ready button is hidden
   - The player is unseated with `unseat_player()`
   - An error message is logged

### Remote Player Seating

When a remote player sits:

1. The server sends a SIT event to all clients
2. `EventManager.SIT_received` signal is emitted
3. `PlayerManager._on_player_sit()` is called with player ID and seat number
4. The method finds the appropriate seat and calls `seat.seat_player(id)`

## Seat Player Logic

The `seat_player(id)` method:

1. Gets the player object from the PlayerManager
2. If the player is already seated elsewhere:
   - Gets the old seat reference
   - Calls `old_seat.unseat_player()`
3. Removes the player from the "pinned players" list
4. Moves the player to the seat's position
5. Updates the player's state to `PLAYER_WAITING`
6. Sets the player's seat reference
7. Updates the seat's sitter reference
8. Disables the seat button

## Unseating Process

### Local Player Unseating

When a local player wants to leave a seat:
1. The player triggers an unseat action (not directly visible in provided code, likely through UI)
2. A network request is sent via `EventManager.unsit_request()`
3. On success, the player's reference is maintained
4. On failure, the player remains seated

### Remote Player Unseating

When a remote player leaves a seat:
1. The server sends an UNSIT event to all clients
2. `EventManager.UNSIT_received` signal is emitted
3. The seat's `unseat_player()` method is called

### Player Disconnection

When any player disconnects:
1. The server sends a LEAVE event
2. `PlayerManager._on_player_leave()` checks if the player was seated
3. If seated, it calls `leaving_player.seat.unseat_player()`

## Unseating Logic

The `unseat_player()` method:
1. If there is a sitter, calls `sitter.unseat()`
2. Sets the sitter reference to null
3. Enables the seat button

The player's `unseat()` method:
1. Sets the player's state to `PLAYER_IDLE`
2. Clears the player's seat reference
3. Moves the player back to the player list via `manager.pin_player(self)`

## Visual Representation

- Each seat has a visual sprite representation
- An interactive button allows players to sit in empty seats
- The button is disabled when a seat is occupied
- The Ready Button appears only when the local player is seated

## Integration with Other Systems

- **Player Manager**: Tracks which players are seated and which are in the lobby
- **Event Manager**: Handles network requests for sitting/unseating
- **Ready Button**: Appears when a player sits down, enabling them to indicate readiness
- **Network Manager**: Ensures all clients have consistent seat assignments

## Limitations and Considerations

- The current implementation supports exactly 4 seats
- Seat positions are fixed at predefined positions around the table
- There is no concept of "reserved" seats - any player can sit in any empty seat
- The system does not handle game-specific seating requirements like team assignments