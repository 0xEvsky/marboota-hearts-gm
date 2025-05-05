# Globals

## Overview

The `Globals` singleton provides a central access point for important game components and references. It serves as a lightweight service locator pattern implementation, allowing various parts of the game to easily access commonly needed objects.

**File:** `scripts/globals.gd`

## Key Responsibilities

1. Providing global access to the PlayerManager instance
2. Providing global access to the local player instance
3. Serving as a central point for component discovery

## Properties

| Property | Type | Description |
|----------|------|-------------|
| player_manager | PlayerManager | Reference to the game's PlayerManager |
| my_player | Player | Reference to the local player instance |

## Usage in the System

The Globals singleton is registered as an autoload in the project settings, making it accessible from any script in the game without needing to use node paths or signals for discovery.

Components that need to be globally accessible register themselves with Globals:
- PlayerManager sets `Globals.player_manager = self` during initialization
- The local player is set with `Globals.my_player = new_player` when created

## Integration with Other Components

- **PlayerManager**: Registers itself with Globals
- **Player**: Local player instance is registered with Globals
- **Seat**: Uses Globals to access player_manager and my_player
- **ReadyButton**: Uses Globals to access my_player

## Access Pattern

The Globals singleton follows a simple property-based access pattern:

```gdscript
# Getting the player manager
var player_manager = Globals.player_manager

# Getting the local player
var my_player = Globals.my_player

# Checking if player is seated
if Globals.my_player.seat != null:
    print("Player is seated!")
```

## Benefits

1. **Simplified Access**: Components can access commonly needed objects without complex node paths
2. **Reduced Coupling**: Components don't need direct references to each other
3. **Centralized Management**: Core component references are managed in one place
4. **Code Readability**: Intent is clear when accessing global components

## Limitations

1. Overuse can lead to tight coupling between components
2. Global state can make testing more difficult
3. No validation when setting properties

## Future Improvements

1. Add type checking for property assignments
2. Implement signal-based notification for when references change
3. Add more core game components as they are developed
4. Consider using a more structured registry pattern for more complex games
