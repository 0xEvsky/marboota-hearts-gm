# ARCHITECTURE.md

## 🏗️ Backend Server Architecture

### Technologies Used:
- **Golang** – high performance, statically typed language ideal for concurrency.
- **WebSockets** – used for real-time, bi-directional communication.

### Design Principles:
- **Modular Structure**: Each component has a clear responsibility.
- **Concurrency via Goroutines**: Efficient parallel player and room handling.
- **Message-driven**: All actions are triggered and propagated via typed messages.

### Component Overview:
- **Main Server (main.go)**: Boots the server and initializes routes.
- **WebSocket Layer (ws/)**: Manages socket connections and routes messages.
- **Room Management (room/)**: Tracks and controls game sessions.
- **Player Management (player/)**: Maintains per-player data and logic.
- **Message Definitions (message/)**: Ensures consistent data formats.
- **Utilities (utils/)**: Supports common operations like ID generation.

### Sequence Diagram (Simplified):
```
Client ──Connect──▶ WS Server
         │               │
         ▼               ▼
   Player Instance ──▶ Room Join/Create
         │               │
         ▼               ▼
     Set Ready     ──▶ Game Start
         │               │
         ▼               ▼
    Broadcast Events ◀─ Game Logic
         │               │
         ▼               ▼
     Disconnect   ──▶ Cleanup
```
