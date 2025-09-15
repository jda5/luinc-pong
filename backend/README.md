# API Reference

The backend provides JSON-based responses and accepts JSON-formatted request bodies.

## GET `/leaderboard`

Fetches the main leaderboard, listing all players sorted by their Elo rating in descending order.

**Response**

Returns an array of player objects, each containing their ID, name, and current Elo rating.

Type: `[]models.LeaderboardRow`

_Example Response_

```json
[
  {
    "id": 1,
    "name": "Alice",
    "eloRating": 1050.5
  },
  {
    "id": 3,
    "name": "Bob",
    "eloRating": 1012.0
  },
  {
    "id": 2,
    "name": "Charlie",
    "eloRating": 980.7
  }
]
```

## GET `/players/:id`

Fetches the detailed profile for a single player, including their stats and their 20 most recent games.

**URL Parameters**

`id` (integer, required): The unique ID of the player to fetch.

**Success Response (200 OK)**

Returns a single PlayerProfile object.
Type: `models.PlayerProfile`

_Example Response_

```json
{
  "id": 1,
  "name": "Alice",
  "eloRating": 1050.5,
  "createdAt": "2023-10-27T10:00:00Z",
  "gamesPlayed": 5,
  "gamesWon": 3,
  "recentGames": [
    {
      "id": 101,
      "winner": {
        "id": 1,
        "name": "Alice"
      },
      "loser": {
        "id": 2,
        "name": "Charlie"
      },
      "winnerScore": 11,
      "loserScore": 5,
      "createdAt": "2023-10-28T14:30:00Z"
    },
    {
      "id": 99,
      "winner": {
        "id": 3,
        "name": "Bob"
      },
      "loser": {
        "id": 1,
        "name": "Alice"
      },
      "winnerScore": null,
      "loserScore": null,
      "createdAt": "2023-10-27T11:00:00Z"
    }
  ]
}
```

## POST `/players`

Creates a new player. The player's name must be unique.

**Request Body**

A JSON object containing the name of the new player.

Type: `models.Name`

_Example Request_

```json
{
  "name": "NewPlayer"
}
```

**Success Response (201 Created)**

Returns a JSON object with the ID of the newly created player.

_Example Response_

```json
{
  "id": 4
}
```

## POST `/games`

Submits the result of a new game. This saves the game to the database and triggers a background task to recalculate the Elo ratings for both players.

**Request Body**

A JSON object detailing the game result. Type: `models.GameResult`

_Example Request (with scores)_

```json
{
  "winnerId": 1,
  "loserId": 2,
  "winnerScore": 11,
  "loserScore": 5
}
```

_Example Request (without scores)_

_The `winnerScore` and `loserScore` fields are optional._

```json
{
  "winnerId": 1,
  "loserId": 2
}
```

**Success Response (201 Created)**

Returns a JSON object with the ID of the newly created game.

_Example Response_

```json
{
  "id": 102
}
```