# Action Handler Documentation

## Action Handler Component

The Action Handler component processes client actions and updates game state accordingly. It validates actions against game rules and current state, then performs the appropriate operations.

### Implementation Files
- `actionHandler.go` (primary implementation)
- `messageHandler.go` (routing to action handlers)

### Primary Responsibilities

1. **Action Validation**:
   - Verify actions are valid for the current game state
   - Check that clients have proper permissions
   - Ensure actions follow game rules

2. **State Modification**:
   - Update game state based on valid actions
   - Transition between game phases when appropriate
   - Update player status and positions

3. **Response Generation**:
   - Send appropriate responses to the acting client
   - Broadcast state changes to other clients

### Key Action Handlers

#### Authentication
```go
func authClient(c *Client, instanceId, userId, userName, iconUrl string) error
```
- Validates client authentication data
- Associates client with an instance
- Sends catch-up data to the client
- Notifies other clients of the new client

#### Seating
```go
func seatClient(c *Client, seatStr string) error
```
- Validates seat request
- Associates client with a player position
- Updates client and player state
- Notifies other clients of the seating change

#### Unseating
```go
func unseatClient(c *Client) error
```
- Validates unseat request
- Disassociates client from player position
- Updates client and player state
- Notifies other clients of the unseating

#### Ready Status
```go
func setReady(c *Client) error
```
- Validates ready request
- Updates player ready status
- Checks if all players are ready
- Starts game if all players are ready

#### Unready Status
```go
func unsetReady(c *Client) error
```
- Validates unready request
- Updates player ready status
- Notifies other clients of the status change

#### Trump Calling
```go
func advanceTrump(c *Client, scoreStr string) error
```
- Validates trump call
- Updates trump state
- Advances turn to next player
- Broadcasts call to other clients

### Action Flow

Each action follows a similar pattern:
1. Validate the action against current state
2. If invalid, return an error
3. Update game state based on the action
4. Notify the client of success
5. Broadcast state changes to other clients
6. Check for game state transitions

### Error Handling

The action handlers perform comprehensive validation:
- Authentication status
- Game state compatibility
- Turn validation
- Action-specific validation

Errors are returned to the client with descriptive messages.

### State Transitions

Action handlers can trigger state transitions:
- `setReady`: May transition from TableWaiting to TableTrumping
- `advanceTrump`: May transition from TableTrumping to TablePlaying

### Game Phase Management

The action handlers coordinate game phase transitions:
1. **Waiting Phase**:
   - Players join, sit, and mark themselves ready
   - When all players are ready, transition to Trumping

2. **Trump Phase**:
   - Players call trump scores or pass
   - When trump is determined, transition to Playing

3. **Playing Phase**:
   - Players play cards in turn
   - Track tricks and scores
   - Determine winners when game ends

### Broadcast Notifications

After successful actions, the handlers broadcast notifications to other clients:
- JOIN: New client has joined
- SIT: Client has sat at a position
- UNSIT: Client has left a position
- READY: Client is ready to play
- UNREADY: Client is no longer ready
- TRUMPCALL: Client has made a trump call

### Example Flow

```
Client sends action
  │
  ▼
Message handler routes to appropriate action handler
  │
  ▼
Action handler validates the action
  │
  ▼
Action handler updates game state
  │
  ▼
Action