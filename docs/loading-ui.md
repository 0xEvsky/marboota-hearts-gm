# Loading UI

## Overview

The `LoadingUI` is responsible for displaying a loading screen while the game is connecting to the server and completing authentication. It provides visual feedback to players during the initialization process.

**File:** `scripts/loading_ui.gd`

## Key Responsibilities

1. Displaying a loading screen during game initialization
2. Providing visual feedback with a spinning loading indicator
3. Hiding once player authentication is complete

## Scene Structure

The LoadingUI is a CanvasLayer with the following elements:

1. **LoadingScreen**: Control node containing all loading UI elements
   - **LoadingBackground**: TextureRect showing background image
   - **LoadingCenter**: CenterContainer for centered content
     - **VBoxContainer**: Container for vertical alignment
       - **LoadingText**: Label displaying "LOADING" text
       - **SpinnerControl**: Control node for spinner positioning
         - **LoadingSpinner**: TextureRect with spinning animation
         - **AnimationPlayer**: Node controlling spinner rotation

## Visual Design

The loading screen features:
- A full-screen background
- Large "LOADING" text in the center
- A continuously spinning icon below the text

## Main Functions

### `_ready()`
Connects to the NetworkManager.AUTH_accepted signal to hide the loading UI once authentication is complete.

## Integration with Other Components

- **NetworkManager**: LoadingUI listens for AUTH_accepted signal to know when to hide
- **Game**: LoadingUI overlays the game scene until authentication is complete

## Loading Process

1. When the game starts, the LoadingUI is visible by default
2. NetworkManager attempts to connect and authenticate with the server
3. Once authentication succeeds, NetworkManager emits AUTH_accepted signal
4. LoadingUI hides in response to this signal, revealing the game

## Example Usage

```gdscript
# The loading UI is automatically created in the game scene
# To manually control visibility:
$LoadingUI.visible = true  # Show loading screen
$LoadingUI.visible = false # Hide loading screen
```

## Future Improvements

1. Add loading progress indication
2. Implement connection status messages
3. Add retry button for failed connections
4. Include animated transitions for smoother UX
