# Networking Protocol Documentation

## Overview

The Marboota Game uses a WebSocket-based client-server communication protocol with JSON messages. This document outlines the structure of messages, request-response patterns, and event handling system.

## Message Format

All messages are JSON objects with an "ACTION" field that determines the message type. The basic structure is:

```json
{
  "ACTION": "ACTION_NAME",
  // Other fields specific to the action
}
```

## Authentication

When the connection is established, the client must authenticate before any other actions:

```json
// Client -> Server
{
  "ACTION": "AUTH", 
  "INSTANCEID": "1234",
  "USERID": "12345", 
  "USERNAME": "Player", 
  "ICONURL": "https://example.com/icon.png"
}

// Server -> Client (Success)
{
  "ACTION": "OK"
}

// Server -> Client (Failure)
{
  "ACTION": "ERROR",
  "MESSAGE": "Error message"
}
```

## Request-Response Pattern

For actions that require server validation, the client uses a request-response pattern:

1. Client sends a request with a unique "REQUESTID"
2. Server responds with the same "REQUESTID" and either "OK" or "ERROR"

```json
// Client -> Server
{
  "ACTION": "SIT",
  "SEAT": "1",
  "REQUESTID": "request-12345-789"
}

// Server -> Client (Success)
{
  "ACTION": "OK",
  "REQUESTID": "request-12345-789"
}

// Server -> Client (Failure)
{
  "ACTION": "ERROR",
  "MESSAGE": "Seat already taken",
  "REQUESTID": "request-12345-789"
}
```

## Game Events

When events occur that all clients need to know about, the server broadcasts notifications:

```json
// Player joins game (Server -> All Clients)
{
  "ACTION": "JOIN",
  "USERID": "12345",
  "USERNAME": "Player",
  "ICONURL": "https://example.com/icon.png"
}

// Player leaves game (Server -> All Clients)
{
  "ACTION": "LEAVE",
  "USERID": "12345"
}

// Player sits at seat (Server -> All Clients)
{
  "ACTION": "SIT",
  "USERID": "12345",
  "SEAT": "1"
}

// Player leaves seat (Server -> All Clients)
{
  "ACTION": "UNSIT",
  "USERID": "12345"
}

// Player is ready (Server -> All Clients)
{
  "ACTION": "READY",
  "USERID": "12345"
}

// Player is no longer ready (Server -> All Clients)
{
  "ACTION": "UNREADY",
  "USERID": "12345"
}
```

## Request Queue System

The `EventManager` implements a request queue system:

1. Client generates a unique `REQUESTID`
2. Request is added to a queue with success/failure callbacks
3. When a response with matching `REQUESTID` is received:
   - The corresponding callbacks are executed
   - The request is removed from the queue

## Available Client Requests

The EventManager exposes these request functions:

```gdscript
# Request to sit at a specific seat
EventManager.sit_request(seat: int) -> Dictionary

# Request to stand up from current seat
EventManager.unsit_request() -> Dictionary

# Indicate player is ready to start game
EventManager.ready_request() -> Dictionary

# Remove ready status
EventManager.unready_request() -> Dictionary
```

Each request can be sent with success and error callbacks:

```gdscript
EventManager.send_request(
    EventManager.sit_request(1),
    func(): # Success handler
        print("Successfully sat down"),
    func(error_msg): # Error handler
        print("Failed to sit down: " + error_msg)
)
```

## Event Signals

The EventManager emits signals when events are received:

```gdscript
# Player joined the game
signal JOIN_received

# Player left the game
signal LEAVE_received

# Player sat at a seat
signal SIT_received

# Player stood up from seat
signal UNSIT_received

# Player is ready
signal READY_received

# Player is no longer ready
signal UNREADY_received
```

## Connection Flow

1. Client connects to WebSocket server
2. Client sends AUTH request
3. Server validates and responds
4. On success, client can start sending game requests

## Error Handling

1. Connection errors: The NetworkManager logs errors and stops processing
2. Authentication errors: The NetworkManager logs error message and stops processing
3. Request errors: Error callbacks handle specific request failures
4. Invalid messages: The EventManager logs errors for unknown message types

## Implementation Details

- NetworkManager handles the WebSocket connection and authentication
- EventManager processes game events and manages the request queue
- All network operations are asynchronous
- The client uses an optimistic UI update approach followed by server confirmation