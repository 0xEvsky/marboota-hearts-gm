# Project Architecture Overview

## Introduction

This document provides a high-level overview of the Marboota Game architecture, explaining how the different components interact with each other and the overall flow of the game system.

## Key Components

The project is structured around these key components:

1. **Network Management**: Handles WebSocket connections, authentication, and message exchange with the server
2. **Event Management**: Processes game events and manages request/response communication
3. **Player Management**: Creates, tracks, and updates player entities
4. **Seating System**: Manages player positions at the game table
5. **Ready System**: Handles player readiness state for game progression

## System Flow

```
                  ┌─────────────────┐
                  │  NetworkManager │
                  └────────┬────────┘
                           │
                           ▼
                  ┌─────────────────┐
                  │  EventManager   │
                  └────────┬────────┘
                           │
             ┌─────────────┴─────────────┐
             │                           │
             ▼                           ▼
   ┌─────────────────┐         ┌─────────────────┐
   │  PlayerManager  │◄────────►      Seat       │
   └─────────────────┘         └─────────────────┘
             │                           │
             ▼                           │
   ┌─────────────────┐                   │
   │     Player      │◄──────────────────┘
   └─────────────────┘
             │
             ▼
   ┌─────────────────┐
   │   ReadyButton   │
   └─────────────────┘
```

## Game Flow

1. The application connects to the WebSocket server through `NetworkManager`
2. User authenticates with username and ID
3. Upon successful authentication, `LoadingUI` disappears
4. `PlayerManager` initializes the local player and receives/manages other connected players
5. Players can sit at available seats around the table
6. Once seated, players can toggle their "ready" state
7. Game logic progresses when all seated players are ready (not fully implemented in current code)

## Global Access

The `Globals` singleton provides convenient access to critical components:
- `player_manager`: Reference to the PlayerManager instance
- `my_player`: Reference to the local player for quick access

## Component Dependencies

| Component | Depends On |
|-----------|------------|
| NetworkManager | WebSocketPeer |
| EventManager | NetworkManager |
| PlayerManager | NetworkManager, EventManager, Globals |
| Player | PlayerManager, Globals |
| Seat | EventManager, Globals |
| ReadyButton | EventManager, Globals |

## Signal Flow

The project uses Godot's signal system extensively for component communication:
1. NetworkManager emits signals when connection state changes
2. EventManager processes network messages and emits appropriate signals
3. PlayerManager and other components subscribe to these signals to update game state

For detailed information about each component, refer to the individual documentation files.
