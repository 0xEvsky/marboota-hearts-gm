# Ready Button

## Overview

The `ReadyButton` manages the player's readiness state. It allows seated players to indicate whether they are ready to begin playing the game, communicating this state to the server.

**File:** `scripts/ready_button.gd`

## Key Responsibilities

1. Managing player readiness state
2. Communicating ready/unready status to the server
3. Visually representing the current readiness state
4. Handling state synchronization with the server

## Properties

The ReadyButton extends Godot's Button class and adds toggle functionality to switch between ready and not ready states.

## Main Functions

### `_ready()`
Initializes the button, hiding it initially since it should only be visible when the local player is seated.

### `_on_toggled(toggled_on)`
Handles button state changes:

When toggled ON:
1. Updates local player state to PLAYER_READY
2. Sends ready request to server
3. If request fails, reverts player state to PLAYER_WAITING

When toggled OFF:
1. Updates local player state to PLAYER_WAITING
2. Sends unready request to server
3. If request fails, reverts player state to PLAYER_READY

## Visual Design

The ReadyButton is designed as a prominent UI element:
- Large green button with white text
- Located at the center of the table
- Text reads "Ready"
- Uses a toggle mode to indicate state

## Integration with Other Components

- **Table**: Parent node that contains the ready button
- **Seat**: Controls button visibility based on local player's seating status
- **Player**: Ready button updates player state
- **EventManager**: Used to send ready/unready requests to the server

## Ready State System

The ready system works as follows:
1. Initially, seated players are in PLAYER_WAITING state and the button is off
2. When a player toggles the button, their state changes to PLAYER_READY
3. The server is notified of this change
4. When all players are ready, the game can proceed to the next phase (not implemented in current code)

## Example Usage

```gdscript
# Accessing the ready button
var ready_button = get_node("Table/ReadyButton")

# Checking if button is toggled
if ready_button.button_pressed:
    print("Player is ready!")
    
# Programmatically toggling the button
ready_button.button_pressed = true
```

## Future Improvements

1. Add visual feedback for other players' ready states
2. Implement timeout for ready state
3. Add animation for state transitions
4. Include counter showing how many players are ready
