# COMPONENTS/MESSAGE.md

## Message Structure (`message/message.go`)

### Purpose:
- Standardize all client-server communication.

### Format:
```go
Message {
  Type string      // e.g., "join", "ready", "start"
  Payload string   // JSON encoded content
}
```

### Key Features:
- Typed messaging simplifies routing.
- Encodes/decodes JSON to/from structs.

### Interactions:
- Used throughout WS, Room, and Player handlers.
