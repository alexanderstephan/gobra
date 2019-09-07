package main

import (
	"container/list"
	gc "github.com/rthornton128/goncurses"
	"log"
	"math/rand"
	"time"
)

const food_char = 'X'
const snake_body = 'O'
const start_size = 5

var snake = list.New()

type Segment struct {
	y, x int
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

var d Direction = East

func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d]
}

func main() {
	// Initialize goncurses
	stdscr, err := gc.Init()

	if err != nil {
		log.Fatal(err)
	}

	// End is required to preserve terminal after execution
	defer gc.End()

	// Randomize pseudo random functions
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

	// Create a rectangle window that is a placeholder for the snake
	var input *gc.Window
	input, err = gc.NewWindow(0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	input.Refresh()

	// Init snake
	InitSnake(stdscr)

	// Setup frame counter
	frame_counter := 0

	// Init food
	newFood := &Food{}

	// Threshold for timeout
	input.Timeout(100)

	// Wait for keyboard input
	input.Keypad(true)

loop:
	for {
		// Clear screen
		stdscr.Refresh()
		stdscr.Erase()

		stdscr.ColorOn(1)

		stdscr.MovePrint(0, 0, "Use vim bindings to move the snake. Press 'q' to exit")
		stdscr.MovePrint(1, 0, "DEBUG:")
		stdscr.MovePrint(2, 0, frame_counter)
		stdscr.MovePrint(3, 0, d)
		stdscr.MovePrint(4, 0, snake.Front().Value.(Segment).y)
		stdscr.MovePrint(4, 3, snake.Front().Value.(Segment).x)

		// Init food position
		if newFood.y == 0 && newFood.x == 0 {
			newFood = &Food{y: rand.Intn(rows), x: rand.Intn(cols)}
		}

		// Detect food collision
		if snake.Front().Value.(Segment).y == newFood.y && snake.Front().Value.(Segment).x == newFood.x {
			newFood = &Food{y: rand.Intn(rows), x: rand.Intn(cols)}
		}

		stdscr.MoveAddChar(newFood.y, newFood.x, food_char)

		// Detect boundaries
		if snake.Front().Value.(Segment).y > rows || snake.Front().Value.(Segment).y < 0 || snake.Front().Value.(Segment).x > cols || snake.Front().Value.(Segment).x < 0 {
			stdscr.MovePrint((rows/2)-1, (cols/2)-4, "GAME OVER")
		}

		// setSnakeDir returns false on exit -> interrupt loop
		if !setSnakeDir(stdscr, input, snake.Front().Value.(Segment).y, snake.Front().Value.(Segment).x) {
			break loop
		}

		// Move snake by one cell in the new direction
		switch d {
		case North:
			MoveSnake(North)
		case South:
			MoveSnake(South)
		case West:
			MoveSnake(West)
		case East:
			MoveSnake(East)
		}

		// Render snake with altered position
		RenderSnake(stdscr)

		stdscr.ColorOff(1)

		// Count frames for debug purposes
		frame_counter++


		// Refresh changes in screen buffer
		stdscr.Refresh()

		// Flush characters that have changed
		gc.Update()
	}
	input.Delete()
}
