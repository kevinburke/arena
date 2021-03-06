package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hoisie/web"
	"github.com/kevinburke/arena/arena"
	"log"
	"net/http"
	"strconv"
	"time"
)

var moveGetter = getMoves

func checkError(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func players(ctx *web.Context) []byte {
	ctx.SetHeader("Content-Type", "application/json", true)
	players, err := arena.GetPlayers()
	checkError(err)
	jsonPlayers, err := json.Marshal(players)
	checkError(err)
	return jsonPlayers
}

type Move struct {
	Player string    `json:"player"`
	Column int       `json:"column"`
	Played time.Time `json:"played"`
}

type Moves struct {
	Moves []*Move `json:"moves"`
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func getMoves(moveId int) []*Move {
	fmt.Println("getting move %d", moveId)
	db := arena.GetConnection()
	// XXX do a join here to get player name
	query := "SELECT fourup_column, player, played FROM fourup_moves WHERE match_id = $1"
	rows, err := db.Query(query, moveId)
	checkError(err)
	var moves []*Move
	for rows.Next() {
		var m Move
		var pId int
		err = rows.Scan(&m.Column, &pId, &m.Played)
		checkError(err)
		player, err := arena.GetPlayerById(pId)
		checkError(err)
		player.SetHref()
		m.Player = player.Href
		moves = append(moves, &m)
	}
	return moves
}

func MovesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["match"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// XXX, 400 error.
		fmt.Println("bad id")
	}
	moves := moveGetter(id)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"moves": moves})
}

func PlayersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	players, err := arena.GetPlayers()
	checkError(err)
	fmt.Fprint(w, Response{"players": players})
}

func InvitationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	invitedPlayerName := mux.Vars(r)["player"]
	fmt.Println(invitedPlayerName)
	fmt.Fprint(w, Response{"player": invitedPlayerName})
}

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func moves(ctx *web.Context, matchId string) []byte {
	db := arena.GetConnection()
	// XXX do a join here to get player name
	query := "SELECT fourup_column, player, played FROM fourup_moves WHERE match_id = $1"
	rows, err := db.Query(query, matchId)
	checkError(err)
	var moves []*Move
	for rows.Next() {
		var m Move
		var pId int
		err = rows.Scan(&m.Column, &pId, &m.Played)
		checkError(err)
		player, err := arena.GetPlayerById(pId)
		checkError(err)
		player.SetHref()
		m.Player = player.Href
		moves = append(moves, &m)
	}
	jsonMoves, err := json.Marshal(moves)
	checkError(err)
	return jsonMoves
}

func DoServer() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/players", PlayersHandler).Methods("GET")
	r.HandleFunc("/games/four-up/matches/{match}/moves", MovesHandler).Methods("GET")
	r.HandleFunc("/players/{player}/invitations", InvitationsHandler).Methods("POST")
	http.Handle("/", r)
	return r
}

func main() {
	router := DoServer()
	log.Fatal(http.ListenAndServe(":8080", router))
}
