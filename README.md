# Battle of Bits

[![Build Status](https://travis-ci.org/battleofbits/arena.png?branch=master)](https://travis-ci.org/battleofbits/arena)

Battle of Bits is a hosted service for competitive board game AI programs. We
host daily tournaments and continual round-robin matches with your robot.

## How it Works

After you sign up on the website, create a profile for your bot. For each bot
your create, you'll need to associate a URL with that bot so that we can talk
to it. Each game a bot plays will need a different URL.

Should a bot only exist for a single type of game? Or should a bot work across
multiple games?

## Webhook API

Each player has an webhook associated webhook for a specifc game. For example,
the player `deepblue` has the following URL associated for Four Up.

    http://example.com/fourup/hook

When it's deepblue's turn, we'll POST to that endpoint with a board state. The
endpoint has 30 seconds to respond with a valid move. If the endpoint takes
longer than 30 seconds, or the returned move isn't valid, the game is
forfeited.  So don't make an invalid move!

### Example Webhook

```js
POST /fourup/hook
Content-Type: application/json

{
  "href": "https://battleofbits.com/games/four-up/matches/1",
  "players": {
    "https://battleofbits.com/players/deepblue": "R",
    "https://battleofbits.com/players/garry": "B"
  },
  "turn": "https://battleofbits.com/players/deepblue",
  "loser": "",
  "winner": "",
  "started": "2013-01-01T23:00:01Z",
  "finished": "",
  "moves": "https://battleofbits.com/games/four-up/matches/1/moves",
  "board": [
    ["","","","","","",""],
    ["","","","","","",""],
    ["","","","","","",""],
    ["","","","","","",""],
    ["","B","","","","",""],
    ["","R","B","","","",""]
  ]
}
```

A valid response is a JSON move.

```js
HTTP/1.0 200 OK
Content-Type: application/json

{
  "column": 2
}
```

## API

Each board game will have a custom media type associated with it, so that
others can utilize the JSON representations, even outside of battleofbits.

### GET /players

```js
{
  "players": [{
    "href": "https://battleofbits.com/players/deepblue",
    "username": "deepblue",
    "name": "Deep Blue"
  },{
    "href": "https://battleofbits.com/players/garry",
    "username": "garry",
    "name": "Garry Kasparov"
  }]
}
```

### POST /players

Create a new player. Username is required

### GET /games

```js
HTTP/1.1 200 OK
Content-Type: application/json

{
  "games": [{
    "name": "Four Up",
    "href": "https://battleofbits.com/games/four-up",
    "matches": "https://battleofbits.com/games/four-up/matches"
  }]
}
```

### GET /games/four-up/matches

Return a list of all present and on-going matches

```js
HTTP/1.1 200 OK
Content-Type: application/json

{
  "matches": [{
    "href": "https://battleofbits.com/games/four-up/matches/1",
    "players": {
      "R": "https://battleofbits.com/players/deepblue",
      "B": "https://battleofbits.com/players/garry"
    },
    "winner": "",
    "started": "2013-01-01T23:00:01Z",
    "finished": "",
    "moves": "https://battleofbits.com/games/four-up/matches/1/moves",
    "board": [
      ["","","","","","",""],
      ["","","","","","",""],
      ["","","","","","",""],
      ["","","","","","",""],
      ["","B","","","","",""],
      ["","R","B","","","",""]
    ]
  }]
}
```

### GET /games/four-up/matches/{id}/moves

```js
{
  "moves": [{
    "player": "https://battleofbits.com/players/deepblue",
    "column": 2,
    "played": "2013-01-01T23:00:01Z",
  }]
}
```

### SUBSCRIBE /games/four-up/matches/{id}/moves

Pass the correct accept header

```
Accept: text/event-stream
```

```
HTTP/1.1 200 OK
Content-Type: text/event-stream

event: patch
data: [
        {
          "op": "replace",
          "path": "/games/four-up/matches/{id}/board"
          "value": [
	      ["","","","","","",""],
	      ["","","","","","",""],
	      ["","","","","","",""],
	      ["","","","","","",""],
	      ["","","","","","",""],
	      ["","","","","","",""]
	    ]
        }
      ]
```

## Developing Battle of Bits

If you wish to work on Battle of Bits itself, you'll first need [Go](http://golang.org)
installed (version 1.1+ is _required_). Make sure you have Go properly
installed, including setting up your
[GOPATH](http://golang.org/doc/code.html#GOPATH).

You'll also need [`lib/pq`](https://github.com/lib/pq)
to compile packer. You can install that with:

```
$ go get -u github.com/lib/pq
```

Next, clone this repository into `$GOPATH/src/github.com/battleofbits/arena` and
then just type `make`. In a few moments, you'll have a working `arena` executable:

```
$ make
...
$ bin/arena
...
```

You can run tests by typing `make test`.

If you make any changes to the code, run `make format` in order to
automatically format the code according to Go standards.
