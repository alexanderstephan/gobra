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
const start_size = 20
const grow_rate = 4

var snake = list.New()
var snake_active bool
var debug_info bool = true

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

func CheckFood(newFood *Food, rows, cols int) bool {
	e := snake.Front()

	for e.Next() != nil {
		if newFood.y == e.Value.(Segment).y && newFood.x == e.Value.(Segment).x {
			return false
		}
		e = e.Next()
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
	if err := gc.InitPair(1, gc.C_GREEN, gc.C_BLACK); err != nil {
		log.Fatal("InitPair failed: ", err)
	}

	if err := gc.InitPair(2, gc.C_RED, gc.C_BLACK); err != nil {
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

		if debug_info == true {
			frame_counter++
			stdscr.MovePrint(1, 0, "DEBUG:")
			stdscr.MovePrint(2, 0, frame_counter)
			stdscr.MovePrint(3, 0, d)
			stdscr.MovePrint(4, 0, snake.Front().Value.(Segment).y)
			stdscr.MovePrint(4, 3, snake.Front().Value.(Segment).x)
			stdscr.MovePrint(5, 0, newFood.y)
			stdscr.MovePrint(5, 3, newFood.x)
			stdscr.MovePrint(6, 0, rows)
			stdscr.MovePrint(6, 3, cols)
			stdscr.MovePrint(7,0, rune(stdscr.MoveInChar(1,20)))

		}

		// Init food position
		if newFood.y == 0 && newFood.x == 0 {
			newFood = &Food{y: rows/2, x: cols/2+15}
		}

		// Detect food collision
		if snake.Front().Value.(Segment).y == newFood.y && snake.Front().Value.(Segment).x == newFood.x {
			newFood = &Food{y: rand.Intn(rows), x: rand.Intn(cols)}
			GrowSnake(grow_rate)
		}


		// setSnakeDir returns false on exit -> interrupt loop
		if !setSnakeDir(input, snake.Front().Value.(Segment).y, snake.Front().Value.(Segment).x) {
			break loop
		}

		// Check if head is element of the body
		e := snake.Front().Next()
		for e != nil {
			if (snake.Front().Value.(Segment).y == e.Value.(Segment).y) && (snake.Front().Value.(Segment).x == e.Value.(Segment).x) {
				stdscr.ColorOn(2)
				stdscr.MovePrint((rows/2)-1, (cols/2)-4, "GAME OVER")
				stdscr.ColorOff(2)
				snake_active = false
			}
			e = e.Next()
		}

		// Detect boundaries
		if (snake.Front().Value.(Segment).y > rows-1) || (snake.Front().Value.(Segment).y < 0) || (snake.Front().Value.(Segment).x > cols-1) || (snake.Front().Value.(Segment).x < 0) {
			stdscr.ColorOn(2)
			stdscr.MovePrint((rows/2)-1, (cols/2)-4, "GAME OVER")
			stdscr.ColorOff(2)
			snake_active = false
		}

		stdscr.ColorOn(1)

		// Render snake with altered position
		if snake_active == true {
			// Move snake by one cell in the new direction
			MoveSnake()
			RenderSnake(stdscr)
		}
		if snake_active == false {
			RenderSnake(stdscr)
		}

		// Draw food
		if stdscr.MoveInChar(newFood.y, newFood.x) == ' ' {
			stdscr.ColorOn(2)
			stdscr.MoveAddChar(newFood.y, newFood.x, food_char)
			stdscr.ColorOff(2)
		} else {
			newFood = &Food{y: rand.Intn(rows), x: rand.Intn(cols)}
		}

		stdscr.ColorOff(1)

		// Refresh changes in screen buffer
		stdscr.Refresh()

		// Flush characters that have changed
		gc.Update()
	}
	input.Delete()
}
