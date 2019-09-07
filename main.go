package main

import (
	"container/list"
	gc "github.com/rthornton128/goncurses"
	"log"
	"math/rand"
	"time"
)

const food_char = 'X'
const snake_alive = 'O'
const snake_dead = '+'
const start_size = 5

var snake = list.New()
var snake_active bool
var debug_info bool = false

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

	if err := gc.InitPair(2, gc.C_BLACK, gc.C_RED); err != nil {
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
	snake_active = true

	// Init food
	newFood := &Food{}

	// Init frame count
	frame_counter := 0

	// Threshold for timeout
	input.Timeout(100)

	// Wait for keyboard input
	input.Keypad(true)

loop:
	for {
		// Clear screen
		stdscr.Refresh()
		stdscr.Erase()


		if ( debug_info == true ) {
			frame_counter++
			stdscr.MovePrint(1, 0, "DEBUG:")
			stdscr.MovePrint(2, 0, frame_counter)
			stdscr.MovePrint(3, 0, d)
			stdscr.MovePrint(4, 0, snake.Front().Value.(Segment).y)
			stdscr.MovePrint(4, 3, snake.Front().Value.(Segment).x)
		}


		// Init food position
		if newFood.y == 0 && newFood.x == 0 {
			newFood = &Food{y: rand.Intn(rows), x: rand.Intn(cols)}
		}


		// Detect food collision
		if snake.Front().Value.(Segment).y == newFood.y && snake.Front().Value.(Segment).x == newFood.x {
			newFood = &Food{y: rand.Intn(rows), x: rand.Intn(cols)}
			GrowSnake(5)
		}
		stdscr.ColorOn(2)

		// Draw food
		stdscr.MoveAddChar(newFood.y, newFood.x, food_char)

		stdscr.ColorOn(2)

		// setSnakeDir returns false on exit -> interrupt loop
		if !setSnakeDir(input, snake.Front().Value.(Segment).y, snake.Front().Value.(Segment).x) {
			break loop
		}

		// Check if head is element of the body
		e := snake.Front().Next()
		for e != nil {
			if (snake.Front().Value.(Segment).y == e.Value.(Segment).y) && (snake.Front().Value.(Segment).x == e.Value.(Segment).x) {
				stdscr.MovePrint((rows/2)-1, (cols/2)-4, "GAME OVER")
				snake_active = false
			}
			e = e.Next()
		}

		// Detect boundaries
		if (snake.Front().Value.(Segment).y > rows) || (snake.Front().Value.(Segment).y < 0) || (snake.Front().Value.(Segment).x > cols) || (snake.Front().Value.(Segment).x < 0) {
			stdscr.MovePrint((rows/2)-1, (cols/2)-4, "GAME OVER")
			snake_active = false
		}

		stdscr.ColorOn(1)

		// Render snake with altered position
		if snake_active == true {
			// Move snake by one cell in the new direction
			MoveSnake()
			RenderSnake(stdscr)
		}
		if snake_active == false  {
			RenderSnake(stdscr)
		}

		stdscr.ColorOff(1)

		// Refresh changes in screen buffer
		stdscr.Refresh()

		// Flush characters that have changed
		gc.Update()
	}
	input.Delete()
}
