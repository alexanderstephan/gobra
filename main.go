package main

import (
	gc "github.com/rthornton128/goncurses"
	"log"
	"math/rand"
	"time"
)

const food_char = 'X'
const snake_body = 'O'
const start_size = 20

type Segment struct {
	y, x int
	next *Segment
}

type Snake struct {
	length int
	start  *Segment
	end    *Segment
}

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

var d Direction = West

func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d]
}

func setSnakeDir(stdscr *gc.Window, input *gc.Window, y, x int) bool {
	// Get screen dimensions
	rows, cols := stdscr.MaxYX()

	// Get input from a dedicated window, otherwise stdscr would be blocked
	k := input.GetChar()

	// Define input handlers with interrupt condition
	switch byte(k) {
	case 'h':
		if x > 0 && d != East {
			d = West
		}
	case 'l':
		if x < cols && d != West {
			d = East
		}
	case 'k':
		if y > 1 && d != South {
			d = North
		}
	case 'j':
		if y < rows && d != North {
			d = South
		}
	case 'q':
		return false
	}
	return true
}

func main() {
	// Initialize goncurses
	stdscr, err := gc.Init()

	if err != nil {
		log.Fatal(err)
	}

	// End is required to preserve terminal after execution
	defer gc.End()

	// Randomize pseudo random fucntions
	rand.Seed(time.Now().Unix())

	// Has the terminal the capability to use color?
	if !gc.HasColors() {
		log.Fatal("Require a color capable terminal")
	}

	// Initalize use of color
	if err := gc.StartColor(); err != nil {
		log.Fatal(err)
	}
	gc.Echo(false)

	// Hide cursor
	gc.Cursor(0)

	// Disable input buffering
	gc.CBreak(true)

	// Set colors
	if err := gc.InitPair(1, gc.C_BLACK, gc.C_GREEN); err != nil {
		log.Fatal("InitPair failed: ", err)
	}

	// Use maximum screen width
	rows, cols := stdscr.MaxYX()

	// Define object dimensions
	y, x := rows/2, cols/2

	// Create a rectangle window that is a placeholder for the snake
	var input *gc.Window
	input, err = gc.NewWindow(0, 0, y, x)
	if err != nil {
		log.Fatal(err)
	}

	input.Refresh()

	// Init snake
	newSnake := &Snake{}
	newSnake.InitSnake(stdscr)

	frame_counter := 0

	// Init food
	var food_y, food_x int
	stdscr.Refresh()

	// Threshold for timeout
	input.Timeout(100)

	// Wait for keyboard input
	input.Keypad(true)

loop:
	for {
		// Clear screen
		stdscr.Refresh()
		stdscr.Erase()

		// Show controls
		stdscr.ColorOn(1)
		stdscr.MovePrint(0, 0, "Use vim bindings to move the snake. Press 'q' to exit")
		stdscr.MovePrint(1, 0, newSnake.length)
		stdscr.MovePrint(3, 0, frame_counter)

		// Init food position
		if food_y == 0 && food_x == 0 {
			food_y = rand.Intn(rows)
			food_x = rand.Intn(cols)
		}

		// Draw food
		stdscr.MoveAddChar(food_y, food_x, food_char)

		// Iterate over list until we get the snakes head
		snake_head := newSnake.start

		for snake_head.next != nil {
			snake_head = snake_head.next
		}

		// Detect food collision
		if snake_head.y == food_y && snake_head.x == food_x {
			food_y = rand.Intn(rows)
			food_x = rand.Intn(cols)
			stdscr.MoveAddChar(food_y, food_x, food_char)
		}

		// Draw new food
		stdscr.MovePrint(2, 0, food_y, food_x)

		// Detect boundaries
		if snake_head.y > rows || snake_head.y < 0 || snake_head.x > cols || snake_head.x < 0 {
			stdscr.MovePrint((rows / 2), cols/2, "GAME OVER")
		}

		// Move snake by one cell in the new direction
		switch d {
		case North:
			newSnake.MoveSnake(North)
		case South:
			newSnake.MoveSnake(South)
		case West:
			newSnake.MoveSnake(West)
		case East:
			newSnake.MoveSnake(East)
		}

		// Cut off unneeded space
		newSnake.CutFront()

		// Count frames for debug purposes
		frame_counter++

		// Render snake with altered position
		newSnake.RenderSnake(stdscr)

		stdscr.ColorOff(1)
		// setSnakeDir returns false on exit -> interrupt loop
		if !setSnakeDir(stdscr, input, y, x) {
			break loop
		}

		// Refresh changes in screen buffer
		stdscr.Refresh()

		// Flush characters that have changed
		gc.Update()
	}
	input.Delete()
}
