# Protocol
The protocol is layered on top of JSON encoding, which is used to carry event information over websockets. The client is expected to comply with the protocol or it will only get errors and desyncs.
Each message must have an `ACTION` key describing the type of event. Furthermore, every client must first authenticate using `AUTH` before the server will accept or send any messages its way.
**ALL** values are expected to be **STRINGS**, even ones that look like numbers.

### REQUESTID
`REQUESTID` is an optional field that can be sent with any request, the server will then echo it back alongside the response (will be `""` if not provided with the request), it can be anything.
It only serves as a way for clients to keep track of which response belongs to which request, which is useful in situations with bad connections (high latency, packetloss...etc).
The `REQUESTID` is expected to be unique for each request, either incremented or random, so the client can actually match request-response pairs even if responses arrive out of order. The server, however, does **NOT** ensure that in any way, it simply echoes back the `REQUESTID` it gets 🤷‍♀️.

### Catch-up
When a client joins after the game has already started for a while, they need to be caught up with the current state of the game; This is done using catch-up messages. The server will quickly send all important past game events one after the other to the client that just joined, things like players joining, cards played, score counts...etc. This is all done using the pre-existing messages for the events themselves so no new handlers are needed in the client.

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

### SETMODE
Requests the server to set the current game mode.
> [!NOTE]
> Valid modes are `"WHIST"` and `"HEARTS"`
```json
{
    "ACTION": "SETMODE",
    "MODE": "WHIST",
    "REQUESTID": "request000123"
}
```

## Server -> client event messages
The server will *- without prompt -* send these messages that contain event updates about game state, other players...etc. Such as notifying all other clients when a client does something (joins, sits..etc).

### JOIN
Whenever a new client authenticates, this message is sent to all other clients in the same instance to inform them of the new client.
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

### PASSCARDS
After the server notifies all clients with `PASSCARDS`, each client will use this to pass the four cards `CARDS` to the appropriate mate and remove them from their hand.
```json
{
    "ACTION": "PASSCARDS",
    "CARDS": "S_13, S_7, S_6, C_10",
}
```


### PLAY
Request to play a card after being prompted by the server with `YOURPLAY`.
```json
{
    "ACTION": "PLAY",
    "CARD": "C_4"
}
```


## Server -> client event messages
The server will *- without prompt -* send these messages that contain event updates about game state, other players...etc. Such as notifying all other clients when a client does something (joins, sits..etc).


### SELECTMODE
Once all players in a table are ready, this is sent to all clients in that instance signaling to select game mode.
```json
{
    "ACTION": "SELECTMODE",
}
```

### GAMESTART
Once the game mode is selected, this is sent to all clients in that instance signaling the game has started.
```json
{
    "ACTION": "GAMESTART"
}
```

### DEAL
Notifies the player of which cards they were dealt randomly, the card names are sent as a comma-seperated string in the `CARDS` field, where the initial letter of the suit name and the power/value of the card are seperated by an underscore `_`, (S_14 = Ace of spades, D_13 = King of diamonds....H_4 = 4 of hearts...etc).
```json
{
    "ACTION": "DEAL",
    "CARDS": "S_14,S_5,S_4,S_2,H_14,H_11,H_10,C_11,C_10,C_9,C_2,D_13,D_2"
}
```

### OTHERDEAL
Notifies a client that all other players at the table have been dealt `COUNT` amount of cards. If the client is a spectator this refers to all 4 players at the table. If the client is one of the players at the table, this refers to the remaining 3 players.
```json
{
    "ACTION": "OTHERDEAL",
    "COUNT": "13"
}
```

### TRUMPSTART
Notifies all clients that trumping has started.
```json
{
    "ACTION": "TRUMPSTART"
}
```

### YOURTRUMPCALL
Notifies the player that it's their turn to call a trump score using `TRUMPCALL`, alongside `MINSCORE` which is the minimum valid score to call (any less will return an error), and `MAXSCORE` which is the maximum valid score to call (any more will return an error).
```json
{
    "ACTION": "YOURTRUMPCALL",
    "MINSCORE": "7",
    "MAXSCORE": "13"
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

### PASSCARDS
Notifies all clients that cards passing phase has started.
```json
{
    "ACTION": "PASSCARDS"
}
```

### RECIEVECARDS
Notifies a client to add the passed cards `CARDS` to their hand.
```json
{
    "ACTION": "RECIEVECARDS",
    "CARDS": "S_13, S_7, S_6, C_10"
}
```

<!-- 
### TRUMPEND
Notifies all players that trump phase is over! Informing them of the selected trump suit and score.
```json
{
    "ACTION": "TRUMPEND",
    "SUIT": "SPADES",
    "SCORE": "9"
}
``` -->

### PLAYSTART
Notifies all players that play has started.
```json
{
    "ACTION": "PLAYSTART"
}
```

### YOURPLAY
Notifies the client that it's their turn to play a card, alongside is the `PLAYABLE` field which holds which cards the client can play.
Just like `DEAL`, the card names are sent as a comma-seperated string in the `CARDS` field, where the initial letter of the suit name and the power/value of the card are seperated by an underscore `_`, (S_14 = Ace of spades, D_13 = King of diamonds....H_4 = 4 of hearts...etc).
```json
{
    "ACTION": "YOURPLAY",
    "PLAYABLE": "S_14,S_5,S_4,S_2"
}
```

### PLAY
Notifies all clients that the user with id `USERID` has played the card `CARD`.
```json
{
    "ACTION": "PLAY",
    "USERID": "11223344",
    "CARD": "H_14"
}
```

### PLAYEND
Notifies all clients that the current hand play has ended, alongside the winner's userid and their penalty.
```json
{
    "ACTION": "PLAYEND",
    "PLAYSCORE": "4",
    "WINNERID": "11223344"
}
```

### TEAMROUNDEND
Notifies all clients that the current hand round (13 plays) has ended. With it are the end scores for the teams for that round `TEAMASCORE` & `TEAMBSCORE`.
```json
{
    "ACTION": "TEAMROUNDEND",
    "TEAMASCORE":  "8",
    "TEAMBSCORE":  "5",
}
```

### TEAMTOTALSCORE
Notifies all clients of the total scores for each team at the moment; `TEAMASCORE` & `TEAMBSCORE`, (sent with `TEAMROUNDEND`).
```json
{
    "ACTION": "TEAMTOTALSCORE",
    "TEAMASCORE":  "15",
    "TEAMBSCORE":  "11",
}
```

### TEAMGAMEEND
Notifies all clients that the game has ended, announcing the USERID of the winning players.
```json
{
    "ACTION": "TEAMGAMEEND",
    "WINNER1ID": "11223344",
    "WINNER2ID": "55667788"
}
```

### FFAROUNDEND
```json
{
    "ACTION": "FFAROUNDEND",
    "0": "1",
    "1": "1",
    "2": "1",
    "3": "10",
}
```

### FFATOTALSCORE
```json
{
    "ACTION": "FFATOTALSCORE",
    "0":  "15",
    "1":  "11",
    "2":  "15",
    "3":  "11",
}
```

### FFAGAMEEND
```json
{
    "ACTION": "FFAGAMEEND",
    "WINNERID": "11223344",
}
```

