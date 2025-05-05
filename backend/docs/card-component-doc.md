# Card Component Documentation

## Card Component

The Card component represents playing cards in the game, with suit, value, and display name properties.

### Implementation Files
- `card.go` (primary definition and methods)

### Structure

```go
type suit int

const (
    Spades suit = iota
    Hearts
    Clubs
    Diamonds
)

type Card struct {
    name  string
    suit  suit
    value int
}
```

### Key Fields

- `name`: String representation of the card (e.g., "S14" for Ace of Spades)
- `suit`: Card suit (Spades, Hearts, Clubs, Diamonds)
- `value`: Card value (2-14, where 14 is Ace)

### Card Naming Convention

Cards are named with a letter representing the suit followed by a number representing the value:

- **Suit Letters**:
  - `S`: Spades
  - `H`: Hearts
  - `C`: Clubs
  - `D`: Diamonds

- **Values**:
  - `2`-`10`: Number cards
  - `11`: Jack
  - `12`: Queen
  - `13`: King
  - `14`: Ace

For example:
- `S14`: Ace of Spades
- `H12`: Queen of Hearts
- `C7`: 7 of Clubs
- `D2`: 2 of Diamonds

### Primary Responsibilities

The Card object:
1. Represents a playing card in the game
2. Provides a string representation for transmission
3. Stores suit and value for game logic

### Key Methods

- `newCard(suit suit, value int)`: Creates a new Card with the given suit and value
- `newDeck()`: Creates a standard 52-card deck

### Deck Generation

The `newDeck()` function creates a complete 52-card deck:
1. Creates an array of 52 Card objects
2. Iterates through all suits (0-3)
3. For each suit, creates cards with values 2-14
4. Returns the complete deck

### Card Sorting

Cards are typically sorted in Player hands by:
1. Suit (primary sort key)
2. Value in descending order (secondary sort key)

This gives a sorted hand with each suit grouped together and cards within each suit ordered from highest to lowest.

### Transmission Format

When sending cards to clients, cards are transmitted as a comma-separated string of card names:

```
"S14,S5,S4,S2,H14,H11,H10,C11,C10,C9,C2,D13,D2"
```

### Usage in Game Logic

Cards are used in several aspects of the game:
1. Dealt to players at the start of the game
2. Played by players during their turns
3. Evaluated for trick-taking and scoring

### Example Flow

```
Game starts
  │
  ▼
Deck of 52 cards is created
  │
  ▼
Deck is shuffled
  │
  ▼
Cards are dealt to players
  │
  ▼
Cards are sorted in player hands
  │
  ▼
Cards are played during gameplay
```

### Interaction with Other Components

- **Table**: The Table creates and deals Card objects
- **Player**: Players hold Card objects in their hands
- **Client**: Card information is transmitted to Clients
