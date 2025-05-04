# User Interface Components Documentation

## Overview

The Marboota Game features several UI layers and components that provide visual feedback, game controls, and information to the player. This document outlines the major UI components, their organization, and functions.

## UI Layers

The game uses Godot's CanvasLayer system to organize UI elements into different layers:

1. **LoadingUI**: Displayed during connection/authentication
2. **DebugUI**: Provides developer information (hidden in release)
3. **BackgroundUI**: Contains the game background image
4. **Game Elements**: Table, seats, and players (not in a separate UI layer)

## Loading UI

### Component Structure
- `LoadingUI` (CanvasLayer)
  - `LoadingScreen` (Control)
    - `LoadingBackground` (TextureRect)
    - `LoadingCenter` (CenterContainer)
      - `VBoxContainer`
        - `LoadingText` (Label)
        - `SpinnerControl` (Control)
          - `LoadingSpinner` (TextureRect with Animation)

### Functionality
- Displays during the connection and authentication process
- Shows an animated spinner to indicate loading activity
- Automatically hides when authentication is successful
- Connected to the `NetworkManager.AUTH_accepted` signal

### Visual Elements
- Background texture
- "LOADING" text with large font size
- Rotating spinner icon

## Debug UI

### Component Structure
- `DebugUI` (CanvasLayer)
  - `DebugPanel` (Control)
    - `Panel`
      - `RichTextLabel`
    - Action Buttons

### Functionality
- Displays debug information about the network connection
- Shows instance ID, username, user ID, icon URL, and authentication status
- Provides test buttons for sitting/unseating
- Hidden by default in the game scene

### Interaction
Three test buttons allow triggering network actions:
- "UNSIT" - Sends an unsit request
- "SIT 1" - Attempts to sit in seat 1
- "SIT 2" - Attempts to sit in seat 2

## Background UI

### Component Structure
- `BackgroundUI` (CanvasLayer)
  - `Background` (TextureRect)

### Functionality
- Provides a consistent background image for the game
- Positioned at layer -1 to appear behind game elements

## Game Table UI

### Component Structure
- `Table` (Node2D)
  - `TableSprite` (Sprite2D)
  - `Seat1` through `Seat4` (Seat scenes)
  - `ReadyButton` (Button)

### Visual Elements
- Circular table sprite in the center
- Four seat positions arranged around the table
- Ready button that appears when a player sits down

### Interaction
- Players can click on empty seats to sit down
- When seated, the Ready button appears and can be toggled

## Seat UI

### Component Structure
- `Seat` (Node2D)
  - `Icon` (Sprite2D)
  - `Button` (Button)

### Visual Elements
- Seat icon showing the position
- Clickable button with the same icon as hitbox

### Interaction
- Click to attempt sitting at the position
- Button is disabled when a seat is occupied

## Ready Button

### Component Structure
- `ReadyButton` (Button) in the Table scene

### Visual Elements
- Text label "Ready"
- Normal state (default button style)
- Pressed state (green background)

### Functionality
- Toggle button that switches between ready/not ready
- Hidden until the local player sits at a seat
- Sends network requests when toggled

## Player Visual Representation

### Component Structure
- `Player` (Node2D)
  - `IconBorder` (Sprite2D)
  - `Icon` (Sprite2D)

### Visual Elements
- Player icon loaded from URL
- Border around the icon (slightly larger scale)

### Positioning
- When in the player list: Positioned along the left side of the screen
- When seated: Positioned at the corresponding seat location

## UI Assets

The game uses several UI assets:
- Font: "SuperBoys-vmW67.ttf"
- Button textures: Various colored rectangle buttons
- Icons: Repeat icon for spinner
- Background image

## UI Styling

The game applies consistent styling through:
- Custom font defined in "basefont.tres"
- Set as the default GUI theme in project settings
- Custom StyleBoxFlat for the Ready button pressed state

## Responsive Layout

- The game's window size is set to 1280x720
- Uses "viewport" stretch mode for consistent scaling
- UI elements use anchors and containers for proper positioning

## UI Flow

1. **Connection**: Loading UI is shown during connection
2. **Lobby**: After connection, shows table with empty seats
3. **Seating**: When a player sits, their icon moves to the seat position
4. **Ready**: When seated, the Ready button appears
5. **Game Start**: When all players are ready, game would transition to gameplay UI
   (Note: The actual gameplay UI is not fully implemented in the provided code)