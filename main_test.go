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
    // Update board for future cases.
    if got {
      placeStone(&c.move)
    }
  }
}

func TestCapture(t *testing.T) {
  // Create groups to capture, check that only those stones are captured.

  // I could have an extra test that creates random groups. Maybe it will find a group that doesn't work properly.

  // Each case should have 2 []Move slices, 1 for groupStones and 1 for surroundingStones.

  game.Turn = 1
  // Single stone
  groupStones := []Move {
    Move{4,4,"black"},
  }
  // L shape
  /*groupStones := []Move {
    Move{4,4},
    Move{5,4},
    Move{6,4},
    Move{6,5},
  }*/

  // FOR EACH GROUP
  // place stones on board
  for _, s := range groupStones {
    placeStone(&s)
  }

  // Create surrounding group and check captures
  surroundingStones := []Move {
    Move{4,3,"white"},
    Move{3,4,"white"},
    Move{5,4,"white"},
    Move{4,5,"white"},
  }

  for _, s := range surroundingStones {
    placeStone(&s)
  }

  // Are captured stones gone? Are surrounding stones still on board? How many points?
  /*
  for i := 0; i < 19; i++ {
    for j := 0; j < 19; j++ {

    }
  }
  */
  for _, s := range groupStones {
    if game.Board[s.X][s.Y] != "" {
      t.Errorf("Move{%d, %d} of groupStones not captured", s.X, s.Y)
    }
  }

  for _, s := range surroundingStones {
    if game.Board[s.X][s.Y] != s.Color {
      t.Errorf("Move{%d, %d} of surroundStones has changed", s.X, s.Y)
    }
  }

  // END FOR EACH GROUP
}