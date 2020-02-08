package main

import (
  "fmt"
  "log"
  "net/http"
  "github.com/gorilla/websocket"
  "encoding/json"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Move)

var upgrader = websocket.Upgrader{}

type Game struct {
  Board [19][19]string `json:"board"`
  Turn int `json:"turn"`
}

var game Game

type Move struct {
  X int `json:"x"`
  Y int `json:"y"`
  Color string `json:"color"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Fatal(err)
  }
  defer ws.Close()
  clients[ws] = true

  for {
    var move Move
    err := ws.ReadJSON(&move)
    if err != nil {
      log.Printf("error: %v", err)
      delete(clients, ws)
      break
    }
    broadcast <- move
    fmt.Println("Received move from player...")
  }
}

// Currently not checking if move is valid, nor saving the game.
func handleMoves() {
  for {
    move := <- broadcast
    // temp
    if game.Turn % 2 == 1 {
      move.Color = "black"
    } else {
      move.Color = "white"
    }
    game.Turn += 1

    game.Board[move.X][move.Y] = move.Color

    fmt.Printf(`Player '%s' made a move at %d-%d.`, move.Color, move.X, move.Y)
    fmt.Printf("\n")
    for client := range clients {
      err := client.WriteJSON(move)
      if err != nil {
        log.Printf("error: %v", err)
        client.Close()
        delete(clients, client)
      }
    }
  }
}

func loadGame(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(game)
}

func main() {
  game.Turn = 1

  fs := http.FileServer(http.Dir("web"))
  http.Handle("/", fs)
  http.HandleFunc("/ws", handleConnections)
  http.HandleFunc("/load", loadGame)

  go handleMoves()

  fmt.Println("Running Go API on localhost:8080...")
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}