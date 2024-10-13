package gameplay

import (
	"container/list"
	"gobra/internal/sound"
	"gobra/internal/tools"
	"os"
	"strconv"
	"time"

	gc "github.com/rthornton128/goncurses"
)

var Run bool

type Config struct {
	Vim       bool
	DebugInfo bool
	NoBounds  bool
	Sound     bool
}

const (
	foodChar   = 'X'
	snakeAlive = 'O'
	snakeDead  = '+'
	snakeHead  = '0'
	startSize  = 5
	growRate   = 3
	scoreMulti = 20
)

var (
	// Objects
	snake   = list.New()
	newFood = Food{}

	// States
	snakeActive bool
	highscore   bool
)

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

var d Direction = East

func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d]
}

// InitSnake draws snake in the screen center pointing east.
func InitSnake(stdscr *gc.Window) {
	screenY, screenX := stdscr.MaxYX()
	for i := 0; i < startSize; i++ {
		snake.PushFront(Segment{y: screenY / 2, x: screenX/2 + i})
	}
}

func Start(cfg *Config) {
	Run = true
	if _, err := os.Stat("/tmp/score"); os.IsNotExist(err) {
		d := []byte("0")
		err = os.WriteFile("/tmp/score", d, 0644)
		tools.Check(err)
	}

	if cfg.Sound {
		sound.InitSound()
	}

	stdscr := SetupGameBoard()
	InitSnake(stdscr) // Create initial snake
	initKeybindings(cfg.Vim)
	input.Timeout(100)              // Threshold for timeout
	tools.Check(input.Keypad(true)) // Wait for keyboard input
	snakeActive = true              // Snake starts alive
	frameCounter := 0               // Init frame count

	stdscr.Refresh()
	time.Sleep(1 * time.Second)

	var scoreLength int
	for Run {
		// Clear the screen.
		stdscr.Refresh()
		stdscr.Erase()

		// Draw box around the screen (for collision detection).
		// TODO: Do we need to redraw everything?
		drawBorder(stdscr)

		// Print debug Infos
		if cfg.DebugInfo {
			frameCounter++
			stdscr.MovePrint(1, 1, "DEBUG:")
			stdscr.MovePrint(2, 1, frameCounter)
			stdscr.MovePrint(3, 1, d)
			stdscr.MovePrint(4, 1, snake.Front().Value.(Segment).y)
			stdscr.MovePrint(4, 4, snake.Front().Value.(Segment).x)
			stdscr.MovePrint(5, 1, newFood.y)
			stdscr.MovePrint(5, 4, newFood.x)
			stdscr.MovePrint(6, 1, screen.rows)
			stdscr.MovePrint(6, 4, screen.cols)
			stdscr.MovePrint(7, 1, rune(stdscr.MoveInChar(0, 0)))
		}

		// setSnakeDir returns false when the user presses q to exit -> interrupt loop.
		if !HandleKeys(input, stdscr, &newFood) {
			break
		}

		// Determine the food position if not set yet.
		initFood(stdscr, &newFood, screen.rows, screen.cols)

		//stdscr.Refresh()

		// Display the snake (alive or dead).
		handleSnake(stdscr, screen.rows, screen.cols)

		// Handle collisions.
		if !handleCollisions(stdscr, &newFood, screen.rows, screen.cols) && cfg.Sound {
			sound.Play(sound.FreqA, 250*time.Millisecond)
		}

		// Check if snake hit boundaries, if desired ports the snake to the other side of the screen.
		boundaryCheck(cfg.NoBounds, screen.rows, screen.cols)

		// Render food symbol.
		printFood(stdscr, &newFood, screen.rows, screen.cols)

		// Overwrite border once again.
		drawBorder(stdscr)
		scoreLength = len(strconv.Itoa(globalScore))

		// Write score to border.
		stdscr.ColorOn(4)
		stdscr.MovePrint(0, (screen.cols/2)-(scoreLength/2), globalScore)
		stdscr.ColorOff(4)
		stdscr.ColorOn(3)
		stdscr.MoveAddChar(0, (screen.cols/2)-(scoreLength/2)-1, '|')

		if scoreLength%2 == 0 {
			stdscr.MoveAddChar(0, (screen.cols/2)+(scoreLength/2), '|')
		} else {
			stdscr.MoveAddChar(0, (screen.cols/2)+(scoreLength/2)+1, '|')
		}

		stdscr.ColorOff(3)
		// Refresh changes in screen buffer.
		stdscr.Refresh()
		// Flush characters that have changed.
		tools.Check(gc.Update())
	}

	gc.End() // Restore previous terminal state.
}

// MoveSnake updates the snake's position based on its current direction.
// It removes the last segment of the snake and inserts a new head segment at the front.
func MoveSnake() {
	// Delete last element of the snake.
	snake.Remove(snake.Back())

	// Read coordinates of the first snake segment.
	headY := snake.Front().Value.(Segment).y
	headX := snake.Front().Value.(Segment).x

	// Increment or decrement last position according to direction.
	switch d {
	case North:
		headY--
	case South:
		headY++
	case West:
		headX--
	case East:
		headX++
	}

	// Insert head with new position.
	snake.PushFront(Segment{headY, headX})
}

// GrowSnake increases the length of the snake by adding segments to its tail.
// It takes an integer parameter 'size' to determine how many segments to add.
func GrowSnake(size int) {
	for i := 0; i < size; i++ {
		tailY := snake.Back().Value.(Segment).y
		tailX := snake.Back().Value.(Segment).x

		// Move segment in the opposite direction.
		switch d {
		case North:
			tailY++
		case South:
			tailY--
		case West:
			tailX++
		case East:
			tailX--
		}

		// Insert segment at back with new position.
		snake.PushBack(Segment{y: tailY, x: tailX})
	}
}

// RenderSnake renders the snake on the provided ncurses window (stdscr).
// It traverses the snake's linked list, drawing each segment based on the snake's state.
// The head of the snake is also drawn.
func RenderSnake(stdscr *gc.Window) {
	// Traverse list and draw every segment to the screen depending on the snake state.
	currentSegment := snake.Front()
	for currentSegment != nil {
		if snakeActive {
			stdscr.MoveAddChar(currentSegment.Value.(Segment).y, currentSegment.Value.(Segment).x, gc.Char(snakeAlive))
		} else {
			stdscr.MoveAddChar(currentSegment.Value.(Segment).y, currentSegment.Value.(Segment).x, gc.Char(snakeDead))
		}
		currentSegment = currentSegment.Next()
	}

	// Attach head.
	if snakeActive {
		stdscr.MoveAddChar(snake.Front().Value.(Segment).y, snake.Front().Value.(Segment).x, snakeHead)
	}
}
