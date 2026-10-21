// Package snake implements the game logic for gitswitch's hidden snake game.
package snake

import "math/rand"

const (
	Width  = 28
	Height = 15
)

// Point identifies one square on the game board.
type Point struct {
	X, Y int
}

// Direction is the direction the snake travels on its next step.
type Direction Point

var (
	Up    = Direction{0, -1}
	Down  = Direction{0, 1}
	Left  = Direction{-1, 0}
	Right = Direction{1, 0}
)

// Game is a single game of snake.
type Game struct {
	Body      []Point
	Food      Point
	Direction Direction
	Score     int
	Over      bool
	Won       bool
	rng       *rand.Rand
}

// NewGame returns a new game with a reproducible food sequence for seed.
func NewGame(seed int64) *Game {
	g := &Game{
		Body:      []Point{{Width / 2, Height / 2}, {Width/2 - 1, Height / 2}, {Width/2 - 2, Height / 2}},
		Direction: Right,
		rng:       rand.New(rand.NewSource(seed)),
	}
	g.placeFood()
	return g
}

// Turn changes direction unless it would make the snake reverse into itself.
func (g *Game) Turn(direction Direction) {
	if g.Over || direction.X+g.Direction.X == 0 && direction.Y+g.Direction.Y == 0 {
		return
	}
	g.Direction = direction
}

// Step advances the game once.
func (g *Game) Step() {
	if g.Over {
		return
	}

	head := g.Body[0]
	next := Point{head.X + g.Direction.X, head.Y + g.Direction.Y}
	if next.X < 0 || next.X >= Width || next.Y < 0 || next.Y >= Height {
		g.Over = true
		return
	}

	ate := next == g.Food
	bodyToCheck := g.Body
	if !ate {
		// The tail leaves its current square on this step, so moving into it is safe.
		bodyToCheck = g.Body[:len(g.Body)-1]
	}
	for _, part := range bodyToCheck {
		if next == part {
			g.Over = true
			return
		}
	}

	g.Body = append([]Point{next}, g.Body...)
	if ate {
		g.Score++
		g.placeFood()
		return
	}
	g.Body = g.Body[:len(g.Body)-1]
}

func (g *Game) placeFood() {
	if len(g.Body) == Width*Height {
		g.Over = true
		g.Won = true
		return
	}
	for {
		food := Point{g.rng.Intn(Width), g.rng.Intn(Height)}
		occupied := false
		for _, part := range g.Body {
			if food == part {
				occupied = true
				break
			}
		}
		if !occupied {
			g.Food = food
			return
		}
	}
}
