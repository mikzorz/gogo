package main

import "testing"

func TestValidMove(t *testing.T) {
  game.Turn = 1

  cases := []struct {
    move Move
    valid bool
  } {
    {Move{X:3,Y:3}, true},
    {Move{X:0,Y:0}, true},
    {Move{X:18,Y:18}, true},
    {Move{X:3,Y:4}, true},
    {Move{X:3,Y:3}, false},
    {Move{X:-1,Y:-1}, false},
    {Move{X:19,Y:19}, false},
  }

  for _, c := range cases {
    got := isMoveValid(c.move)
    if got != c.valid {
      t.Errorf("isMoveValid(%v) == %t, want %t", c.move, got, c.valid)
    }
    if got {
      placeStone(&c.move)
    }
  }
}