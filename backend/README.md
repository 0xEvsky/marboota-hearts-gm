# How to run
- Install Go v1.24
- `cd backend`
- `go run .` to run (as dev, for debugging,etc...)
- `go build .` to compile into a binary for your system

# [Server model](./MODEL.md)
Doc explaining the structural model of the various parts making up the server.

# [Protocol specs](./PROTOCOL.md)
Specifications & details for the protocol used on top of JSON encoding used on top of websockets for sharing game events between the clients.