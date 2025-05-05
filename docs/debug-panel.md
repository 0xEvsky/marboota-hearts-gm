# Debug Panel

## Overview

The `DebugPanel` provides a development and testing interface for the game. It displays connection information and allows developers to trigger various game actions for testing purposes.

**File:** `scripts/debug_panel.gd`

## Key Responsibilities

1. Displaying network connection information
2. Providing test buttons for triggering game actions
3. Showing real-time updates of player state

## Scene Structure

The DebugPanel is a Control node with the following elements:

1. **Panel**: Background panel
   - **RichTextLabel**: Text display area for debug information
2. **Button**: "UNSIT" button to test unseating
3. **Button2**: "SIT 1" button to test sitting at seat 1
4. **Button3**: "SIT 2" button to test sitting at seat 2

## Main Functions

### `_process(_delta)`
Updates the debug information in the RichTextLabel every frame, displaying:
- Instance ID
- Username
- User ID
- Icon URL
- Authentication status

### `_on_button_button_up()`
Sends an UNSIT request when the UNSIT button is pressed.

### `_on_button_2_button_up()`
Sends a SIT request for seat 1 when the SIT 1 button is pressed.

### `_on_button_3_button_up()`
Sends a SIT request for seat 2 when the SIT 2 button is pressed.

## Debug Information

The panel displays the following information in real-time:
- **instance id**: The game instance identifier
- **username**: Current player's display name
- **user id**: Current player's unique identifier
- **icon url**: URL to player's avatar image
- **authenticated**: Whether player is authenticated with server

## Integration with Other Components

- **NetworkManager**: DebugPanel reads connection information from NetworkManager
- **EventManager**: DebugPanel uses EventManager to send test requests
- **Game**: DebugPanel is included in a DebugUI layer that can be toggled

## Usage During Development

The DebugPanel is hidden by default in the game scene. To use it:

1. Set the DebugUI CanvasLayer visibility to true in the inspector
2. Run the game
3. Observe connection information
4. Use the buttons to test seating functionality

## Example Usage

```gdscript
# Accessing the debug panel
var debug_panel = get_node("DebugUI/DebugPanel")

# Showing debug panel
$DebugUI.visible = true

# Hiding debug panel
$DebugUI.visible = false
```

## Future Improvements

1. Add more test buttons for additional game actions
2. Implement a toggle key to show/hide the panel
3. Add logging functionality
4. Include network traffic monitoring
