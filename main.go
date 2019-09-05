package main

import (
	gc "github.com/rthornton128/goncurses"
	"log"
	"math/rand"
	"time"
)

const food_char = 'X'
const snake_body = '#'
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

type Food struct {
	y, x int
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
	rows, cols := stdscr.MaxYX()
	k := input.GetChar()

	switch byte(k) {
	case 'q':
		return false
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

		if food_y == 0 && food_x == 0 {
			food_y = rand.Intn(rows)
			food_x = rand.Intn(cols)
		}
		stdscr.MoveAddChar(food_y, food_x, food_char)
		snake_tail := newSnake.start

		for snake_tail.next != nil {
			snake_tail = snake_tail.next
		}

		if snake_tail.y == food_y && snake_tail.x == food_x {
			food_y = rand.Intn(rows)
			food_x = rand.Intn(cols)
			stdscr.MoveAddChar(food_y, food_x, food_char)
		}

		stdscr.MovePrint(2, 0, food_y, food_x)

		if snake_tail.y > rows || snake_tail.y < 0 || snake_tail.x > cols || snake_tail.x < 0 {
			stdscr.MovePrint((rows / 2), cols/2, "GAME OVER")
		}

		stdscr.ColorOff(1)

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
		newSnake.CutFront()

		frame_counter++

		// Render snake with altered position
		newSnake.RenderSnake(stdscr)

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
