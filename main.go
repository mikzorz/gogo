package main

import (
  "fmt"
  "log"
  "net/http"
  "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Move)

var upgrader = websocket.Upgrader{}

type Game struct {
  Board [19][19]string
}

var currentMove = 1

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

func handleMoves() {
  for {
    move := <- broadcast
    // temp
    if currentMove % 2 == 1 {
      move.Color = "black"
    } else {
      move.Color = "white"
    }
    currentMove += 1

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

func main() {
  fs := http.FileServer(http.Dir("web"))
  http.Handle("/", fs)
  http.HandleFunc("/ws", handleConnections)

  go handleMoves()

  fmt.Println("Running Go API on localhost:8080...")
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}