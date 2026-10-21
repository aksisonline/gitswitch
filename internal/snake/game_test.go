package snake

import "testing"

func TestStepMovesSnake(t *testing.T) {
	g := NewGame(1)
	before := g.Body[0]

	g.Step()

	if got := g.Body[0]; got != (Point{before.X + 1, before.Y}) {
		t.Fatalf("head = %#v, want one square right of %#v", got, before)
	}
	if len(g.Body) != 3 {
		t.Fatalf("length = %d, want 3", len(g.Body))
	}
}

func TestStepEatingFoodGrowsSnake(t *testing.T) {
	g := NewGame(1)
	g.Food = Point{Width/2 + 1, Height / 2}

	g.Step()

	if g.Score != 1 {
		t.Fatalf("score = %d, want 1", g.Score)
	}
	if len(g.Body) != 4 {
		t.Fatalf("length = %d, want 4", len(g.Body))
	}
}

func TestTurnDoesNotAllowReverse(t *testing.T) {
	g := NewGame(1)

	g.Turn(Left)

	if g.Direction != Right {
		t.Fatalf("direction = %#v, want %#v", g.Direction, Right)
	}
}

func TestStepEndsGameAtWall(t *testing.T) {
	g := NewGame(1)
	g.Body = []Point{{Width - 1, 0}, {Width - 2, 0}, {Width - 3, 0}}

	g.Step()

	if !g.Over {
		t.Fatal("game should end when the snake reaches the wall")
	}
}
