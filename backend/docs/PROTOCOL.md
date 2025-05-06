# Protocol
The protocol is layered on top of JSON encoding, which is used to carry event information over websockets. The client is expected to comply with the protocol or it will only get errors and desyncs.
Each message must have an `ACTION` key describing the type of event. Furthermore, every client must first authenticate using `AUTH` before the server will accept or send any messages its way.
**ALL** values are expected to be **STRINGS**, even ones that look like numbers.

### REQUESTID
`REQUESTID` is an optional field that can be sent with any request, the server will then echo it back alongside the response (will be `""` if not provided with the request), it can be anything.
It only serves as a way for clients to keep track of which response belongs to which request, which is useful in situations with bad connections (high latency, packetloss...etc).
The `REQUESTID` is expected to be unique for each request, either incremented or random, so the client can actually match request-response pairs even if responses arrive out of order. The server, however, does **NOT** ensure that in any way, it simply echoes back the `REQUESTID` it gets 🤷‍♀️.

## Server responses
The server responds with either of these to client requests.

### OK
Everything is A-OK 👍.
```json
{
    "ACTION":"OK",
    "REQUESTID": "request000123"
}
```

### ERROR
Error with message.
```json
{
    "ACTION":"ERROR",
    "MESSAGE":"your error message will be here",
    "REQUESTID": "request000123"
}
```

# Lobby events

## Client -> server requests
These are messages clients can send to the server to request a specific action.

### AUTH
Registers the client & its details in the server, required before any further communication is established (the server will **NOT** send any event updates and will only reply with errors if not authenticated).
```json
{
    "ACTION": "AUTH",
    "INSTANCEID": "1234",
    "USERID": "11223344",
    "USERNAME": "Psycho",
    "ICONURL": "discord.com/avatar/11223344.png",
    "REQUESTID": "request000123"
}
```

### SIT
Requests to sit at the provided seat (1-4) for players. Returns an error if a player seat was requested and it was taken or invalid.
```json
{
    "ACTION":"SIT",
    "SEAT":"1",
    "REQUESTID": "request000123"
}
```

### UNSIT
Requests to unsit the client from the game table and back to spectating. Returns an error if the player was not seated to begin with or the game has already started.
```json
{
    "ACTION":"UNSIT",
    "REQUESTID": "request000123"
}
```

### READY
Requests to set the client at the game table as ready. Returns an error if the player is not seated or was already ready.
> [!NOTE]
> Once all players in a table are ready, the game will trigger the PLAYSTART sequence
```json
{
    "ACTION":"READY",
    "REQUESTID": "request000123"
}
```

### UNREADY
Requests to set the client at the game table as NOT ready. Returns an error if the player is not seated, wasn't ready or the game has already started.
```json
{
    "ACTION":"UNREADY",
    "REQUESTID": "request000123"
}
```

## Server -> client event messages
The server will *- without prompt -* send these messages that contain event updates about game state, other players...etc. Such as notifying all other clients when a client does something (joins, sits..etc).

### JOIN
Whenever a new client authenticates, this message is sent to all other clients in the same instance to inform them of the new client.

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
This is sent to all other clients in an instance when a client is successfully seated, alongside its information.
> [!NOTE]
> This is also catch-up sent just like `JOIN`
```json
{
    "ACTION": "SIT",
    "SEAT": "1",
    "USERID": "11223344"
}
```

### UNSIT
This is sent to all other clients in an instance when a client is Unseated, alongside its information.
```json
{
    "ACTION": "UNSIT",
    "USERID": "11223344"
}
```

### READY
This is sent to all other clients in an instance when a client is successfully set as ready, alongside its information.
> [!NOTE]
> This is also catch-up sent just like `JOIN`
```json
{
    "ACTION": "READY",
    "USERID": "11223344"
}
```

### UNREADY
This is sent to all other clients in an instance when a client is set as unready, alongside its information.
```json
{
    "ACTION": "UNREADY",
    "USERID": "11223344"
}
```

# Game events

## Client -> server requests

### TRUMPCALL
After the server notifies the client with `YOURTRUMPCALL`, the client should use this to submit a score as a trump call (تسمية), returns an error if the call is invalid (invalid score, invalid turn...etc).
```json
{
    "ACTION": "TRUMPCALL",
    "SCORE": "8"
}
```
When called with the `SCORE` value being `PASS`, the call is considered a pass in game terms.
```json
{
    "ACTION": "TRUMPCALL",
    "SCORE": "PASS"
}
```

### TRUMPSUIT
After the server notifies the client with `YOURTRUMPSUIT` when trumping is almost over and it's the highest caller, the client should use this to submit the suit for the trump.
```json
{
    "ACTION": "TRUMPSUIT",
    "SCORE": "DIAMONDS"
}
```

## Server -> client event messages
The server will *- without prompt -* send these messages that contain event updates about game state, other players...etc. Such as notifying all other clients when a client does something (joins, sits..etc).

### DEAL
Notifies the player of which cards they were dealt randomly, the card names are sent as a comma-seperated string in the `CARDS` field, where the first letter is the initial letter of the suit name, and the number following it is the power/value of the card, (14: Ace, 13: King....3: 3, 2: 2).
```json
{
    "ACTION": "DEAL",
    "CARDS": "S14,S5,S4,S2,H14,H11,H10,C11,C10,C9,C2,D13,D2"
}
```

### TRUMPSTART
Once all players in a table are ready, this is sent to all clients in that instance signaling the game has started at the TRUMPING state.
```json
{
    "ACTION": "TRUMPSTART"
}
```

### YOURTRUMPCALL
Notifies the player that it's their turn to call a trump score using `TRUMPCALL`, alongside `MINSCORE` which is the minimum valid score to call (any less will return an error).
```json
{
    "ACTION": "YOURTRUMPCALL",
    "MINSCORE": "7"
}
```

### TRUMPCALL
Notifies a player that another player has made a successful trump call, along with the `USERID` of said player, and the `SCORE` called (being `PASS` if it's a pass).
```json
{
    "ACTION": "TRUMPCALL",
    "USERID": "11223344",
    "SCORE": "8"
}
```

### YOURTRUMPSUIT
Notifies the player with the highest call that they should call their trump suit using `TRUMPSUIT`, alongside it is the highest called `SCORE`.
```json
{
    "ACTION": "YOURTRUMPSUIT",
    "SCORE": "10"
}
```

### TRUMPEND
Notifies all players that trump phase is over! Informing them of the selected trump suit and score.
```json
{
    "ACTION": "TRUMPEND",
    "SUIT": "SPADES",
    "SCORE": "9"
}
```