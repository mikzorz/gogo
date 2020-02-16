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

// Currently not checking if move is valid, nor saving the game to storage.
func handleMoves() {
  for {
    move := <- broadcast
    if isMoveValid(move) {
      // For sake of inbounds(), pass value THEN convert to pointer, or use multiple inbounds funcs?
      placeStone(&move)
      fmt.Printf(`Player '%s' made a move at %d-%d.`, move.Color, move.X, move.Y)
      fmt.Printf("\n")
    }

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

// Turn into struct methods?
func isMoveValid(m Move) bool {
  // doesn't check for illegal moves or whose turn it is
  if !inbounds(m) {
    return false
  }
  if game.Board[m.X][m.Y] != "" {
    return false
  }
  return true
}

func placeStone(m *Move) {
  // simple, doesn't check for captures
  if m.Color == "" && game.Turn % 2 == 1 {
    m.Color = "black"
  } else {
    m.Color = "white"
  }
  game.Board[m.X][m.Y] = m.Color

  checkCaptures(m)

  game.Turn += 1
}

func checkCaptures(m *Move) {
  // If neighbour is inbounds
  //   If colour is opposite
  //     If not in list
  //       Add to list, move to neighbours.
  //   If neighbour colour is "", stop. (Either no group to capture or group not captured)
  // Repeat, but check for SAME colour instead.
  // When finished, group is considered captured.
  // Hopefully, this is enough for any type of shape.
  // If group is captured, remove it before checking the next direction.

  // group should be pointers?
  group := []Move

  // Check up
  neighbour := Move{m.X,m.Y-1}
  if inbounds(neighbour) {
    if isOtherColor(neighbour, m.Color) {
      group = append(group, neighbour)
    } else {
      // No neighbour
      // Break loop
    }
  }

  // If this is a recursive function, will it end automatically after all stones are found?
  // Could also have a second list, add neighbour in above section, next while that list is not empty...
  // Get all stones of group
  // Check all four directions
  if inbounds(neighbour) {
    if !isOtherColor(neighbour, m.Color) {
      if !isMoveInGroup(neighbour, group) {
        group = append(group, neighbour)
      }
    } else {
      // No neighbour
      // Break loop
    }
  }

  // All stones scouted.
  // Return group.

}

func inbounds(m Move) bool {
  if m.X < 0 || m.X >= 19 || m.Y < 0 || m.Y >= 19 {
    return false
  }
  return true
}

func isOtherColor(target Move, c string) bool {
  // c = color to check against. If c = black, return true if target.Color == white.
  tc := target.Color
  if tc == "" {
    tc = game.Board[target.X][target.Y].Color
  }

  if tc != "" {
    if tc != c {
      return true
    } else {
      return false
    }
  } else {
    return false
  }

  // If, for some reason, tc is a weird value.
  return false
}

func isMoveInGroup(m Move, s []Move) bool {
  for i := range s {
    if s[i].X == m.X && s[i].Y == m.Y {
      return true
    }
  }
  return false
}

func loadGame(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(game)
}

func resetGame(w http.ResponseWriter, r *http.Request) {
  for i := range game.Board {
    for j := range game.Board[i] {
      game.Board[i][j] = ""
    }
  }
  game.Turn = 1
  fmt.Println("Game reset")

  json.NewEncoder(w).Encode(game)
}

func main() {
  game.Turn = 1

  fs := http.FileServer(http.Dir("web"))
  http.Handle("/", fs)
  http.HandleFunc("/ws", handleConnections)
  http.HandleFunc("/load", loadGame)
  http.HandleFunc("/reset", resetGame)

  go handleMoves()

  fmt.Println("Running Go API on localhost:8080...")
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}