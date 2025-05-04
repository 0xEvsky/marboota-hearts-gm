# Game Flow Documentation

## Overview

The game flow in the Marboota Game follows a state-based architecture that progresses through several distinct phases. This document outlines how the game flows from connection to gameplay.

## Game States

The game uses an enumeration in the PlayerManager to track player states:

```gdscript
enum {PLAYER_UNAVAILABLE, PLAYER_IDLE, PLAYER_WAITING, PLAYER_READY, PLAYER_TRUMPING, PLAYER_PLAYING}
```

### State Descriptions

- **PLAYER_UNAVAILABLE**: Player exists but cannot participate (initial state before connection)
- **PLAYER_IDLE**: Player is connected but not seated at a table
- **PLAYER_WAITING**: Player is seated but not ready to play
- **PLAYER_READY**: Player is seated and has indicated they're ready to play
- **PLAYER_TRUMPING**: Player is in the trump selection phase (game specific)
- **PLAYER_PLAYING**: Player is actively playing in a game round

## Connection Flow

1. **Loading Screen**: When the game starts, the LoadingUI is shown
2. **Authentication**: The NetworkManager attempts to authenticate with the server
3. **Connection Established**: Upon successful authentication, the LoadingUI is hidden
4. **Player Creation**: A player object is created for the local user via `_init_my_player()`
5. **Lobby Entry**: Player enters the "lobby area" in PLAYER_IDLE state

## Game Preparation Flow

1. **Seat Selection**: Player clicks on an available seat
   - Local state changes to PLAYER_WAITING
   - A request is sent to the server to confirm seat selection
   - On success, the player remains seated
   - On failure, the player returns to PLAYER_IDLE state

2. **Ready State**: After sitting, player can press the Ready button
   - Local state changes to PLAYER_READY
   - A request is sent to the server to confirm ready status
   - On success, the player remains in PLAYER_READY state
   - On failure, the player returns to PLAYER_WAITING state

3. **Game Start**: When all seated players are ready (server-determined)
   - Game would transition to gameplay phases (PLAYER_TRUMPING, PLAYER_PLAYING)
   - *Note: The current implementation doesn't contain full gameplay logic*

## Unseating Flow

1. **Manual Unseating**: Player can choose to leave a seat
   - Player clicks a UI element or button to unseat
   - A request is sent to the server
   - On success, player returns to PLAYER_IDLE state
   - Player is moved back to the player list

2. **Disconnection Unseating**: When a player disconnects
   - The server sends a LEAVE notification
   - Other clients unseat the disconnected player
   - The player object is removed from the game

## State Transitions

The following diagram illustrates the primary state transitions:

```
[CONNECTION] → PLAYER_IDLE → PLAYER_WAITING → PLAYER_READY → [GAMEPLAY]
                    ↑               ↑
                    |               |
                    └───────────────┘
                     (unseating)
```

## Network Synchronization

State changes follow a pattern:
1. Local state change (optimistic UI update)
2. Server request to confirm change
3. On success: keep the new state
4. On failure: revert to previous state

This approach gives responsive feedback while maintaining synchronization with the server.

## Key Interactions

- **PlayerManager**: Coordinates player transitions between states
- **EventManager**: Handles network events related to game flow
- **Seat**: Manages player seating and unseating
- **ReadyButton**: Controls player readiness state

## Future Gameplay Flow

The current implementation focuses on the connection, seating, and ready phases. The actual card gameplay states (PLAYER_TRUMPING, PLAYER_PLAYING) are defined but not fully implemented in the provided code.