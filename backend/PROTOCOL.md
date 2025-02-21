# Protocol
The protocol is layered on top of JSON encoding, which is used to carry event information over websockets. The client is expected to comply with the protocol or it will only get errors and desyncs.
Each message must have an `ACTION` key describing the type of event. Furthermore, every client must first authenticate using `AUTH` before the server will accept or send any messages its way.
## Client -> server messages
### AUTH
Registers the client & its details in the server, required before any further communication is established.
```json
{
    "ACTION": "AUTH",
    "INSTANCEID": "1234",
    "USERID": "11223344",
    "USERNAME": "Psycho",
    "ICONURL": "discord.com/avatar/11223344"
}
```
### SIT
Requests to sit at the game table, returns an error if the table is full, otherwise seats the client at the first seat available (1-4).
```json
{
    "ACTION":"SIT"
}
```
### UNSIT
Requests to unsit from the game table, returning to the spectating benches.
```json
{
    "ACTION":"UNSIT"
}
```
### SWITCH
Requests to switch teams at the game table, other team must have a spot open. Must be seated first.
```json
{
    "ACTION":"SWITCH"
}
```
## Server -> client messages
### OK
Everything is A-OK 👍.
```json
{
    "ACTION":"OK"
}
```
### ERROR
Error with message.
```json
{
    "ACTION":"ERROR",
    "MESSAGE":"Your error message will be here"
}
```
### JOIN
Whenever a new client authenticates, this message is sent to all other clients in the same instance that were already connected to inform them of the new client.
> [!NOTE]
> The server will also send multiple JOIN messages to the new client, informing it of the members that were already connected before (catch-up).
```json
{
    "ACTION": "JOIN",
    "USERID": "55667788",
    "USERNAME": "Psycho",
    "ICONURL": "discord.com/avatar/55667788"
}
```
### LEAVE
Sent to all clients in the instance when a client disconnects announcing its user ID
```json
{
    "ACTION":"LEAVE",
    "USERID":"11223344"
}
```
### SIT
This is sent to all other clients in an instance when a client is successfully seated alongside its information.
> [!NOTE]
> This is catch-up sent just like `JOIN`
```json
{
    "ACTION": "SIT",
    "SEAT": "1",
    "USERID": "11223344"
}
```
### UNSIT
This is sent to all other clients in an instance when a client is successfully unseated alongside its information.
```json
{
    "ACTION": "UNSIT",
    "USERID": "11223344"
}
```