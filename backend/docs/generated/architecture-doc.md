# Server Architecture Documentation

## Overview

The Marboota backend is a WebSocket-based server for a card game, with a carefully designed architecture that separates concerns between connection management, game state, and player interactions. This document explains the overall architecture and how the components fit together.

## Architecture Diagram

```
┌────────────────────────────────────────────────────┐
│                      Server                        │
│  ┌──────────────┐         ┌───────────────────┐    │
│  │  WebSocket   │◄────────┤  Message Handler  │    │
│  │  Connection  │         └───────┬───────────┘    │
│  └──────┬───────┘                 │                │
│         │                         │                │
│         ▼                         ▼                │
│  ┌──────────────┐         ┌───────────────────┐    │
│  │  Client Map  │◄────────┤  Action Handler   │    │
│  └──────┬───────┘         └───────┬───────────┘    │
│         │                         │                │
│         ▼                         ▼                │
│  ┌──────────────┐         ┌───────────────────┐    │
│  │   Instance   │◄────────┤  Game State Logic │    │
│  │     Map      │         └───────────────────┘    │
│  └──────────────┘                                  │
└────────────────────────────────────────────────────┘

                        │
                        ▼

┌──────────────────────────────────────────────────────┐
│                     Instance                         │
│  ┌─────────────┐              ┌─────────────────┐    │
│  │  Client Map │◄─────────────┤  Table (Game)   │    │
│  └─────────────┘              └────────┬────────┘    │
│                                        │             │
│                                        ▼             │
│                               ┌─────────────────┐    │
│                               │  Player Array   │    │
│                               └─────────────────┘    │
└──────────────────────────────────────────────────────┘
```

## Core Architecture Principles

1. **Hierarchical Structure**: The server follows a clear hierarchical structure from Server → Instance → Table → Players.

2. **Separation of Concerns**: 
   - Connection management (WebSocket handling) 
   - Game logic (Table, Players, Cards)
   - Message routing (Action handlers)

3. **Message-Driven Communication**: All interactions are done through JSON-encoded messages with a well-defined protocol.

4. **State Management**: Each component maintains its own state and is responsible for its own data integrity.

5. **Concurrency Control**: Mutex locks are used to handle concurrent access to shared resources.

## Component Relationships

### Server to Instance Relationship
- The Server maintains a map of all running Instances
- Each Instance has a unique ID
- The Server routes WebSocket messages to the appropriate Instance

### Instance to Table Relationship
- Each Instance contains exactly one Table
- The Table within an Instance manages the game state
- The Instance provides a bridge between connected Clients and the Table

### Table to Player Relationship
- A Table has a fixed array of 4 Player positions
- Each Player position may be filled or empty
- The Table manages turns, game phases, and game rules

### Client to Player Relationship
- A Client may be associated with a Player when seated
- Not all Clients are associated with Players (spectators)
- The Client provides the communication channel for the Player

## Data Flow

1. **Client Connection**:
   - WebSocket connection established
   - Server creates Client object
   - Client authenticates and is added to an Instance

2. **Game Setup**:
   - Client sends sit request
   - Server assigns Client to a Player position
   - Client sends ready status
   - Server initiates game when all Players are ready

3. **Game Progress**:
   - Server deals cards to Players
   - Table manages turn order and game rules
   - Clients send actions (trump calls, card plays)
   - Server validates actions and updates game state
   - Server broadcasts state changes to all Clients

4. **Client Disconnection**:
   - WebSocket connection closes
   - Server removes Client from Instance
   - Server updates game state if needed
   - Server broadcasts Client departure

## Thread Safety

The server employs mutex locks at several levels:
- Server-level mutex for access to connection and instance maps
- Instance-level mutex for access to client map and table state
- Client-level mutex for writing to WebSocket connection

These locks ensure that concurrent operations (like multiple clients acting simultaneously) don't corrupt the game state.

## System Boundaries

- **Input Boundary**: WebSocket endpoint (`/ws`)
- **Output Boundary**: JSON messages to clients
- **Persistence Boundary**: Currently in-memory only
- **External Interfaces**: None (self-contained system)

## Future Architecture Considerations

- Database persistence for game state
- Authentication service integration
- Metrics and monitoring
- Horizontal scaling across multiple servers
