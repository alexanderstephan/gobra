package main

import (
	gc "github.com/rthornton128/goncurses"
	"log"
	"math/rand"
	"time"
)

const food_char = 'X'

const snake_ascii = "####"

type Direction int

type Food struct {
	*gc.Window
	eaten bool
	y, x  int
}

type Snake struct {
	*gc.Window
	alive     bool
	y, x      int
	segements int
}

type Board struct {
	Snake
	Food
	max_y, max_x int
}

const (
	North Direction = iota
	East
	South
	West
)

var d Direction = North

func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d]
}

func spawnFood(stdscr *gc.Window) *Food {
	val1, val2 := stdscr.MaxYX()
	y := rand.Intn(val1)
	x := rand.Intn(val2)
	stdscr.MoveAddChar(y, x, food_char)
	return &Food{eaten: false, y: y, x: x}
}

func handleInput(stdscr *gc.Window, snake *gc.Window) bool {
	rows, cols := stdscr.MaxYX()
	y, x := snake.YX()
	k := snake.GetChar()

	switch byte(k) {
	case 'q':
		return false
	case 'h':
		if x > 0 {
			x--
		}
	case 'l':
		if x < cols {
			x++
		}
	case 'k':
		if y > 1 {
			y--
		}
	case 'j':
		if y < rows {
			y++
		}
	default:
		return false
	}
	snake.Erase()
	snake.Refresh()
	snake.MoveWindow(y, x)
	snake.Print(snake_ascii)
	snake.Refresh()
	return true
}

func main() {
	// Initialize goncurses
	// End is required to preserve terminal after execution
	stdscr, err := gc.Init()
	if err != nil {
		log.Fatal(err)
	}
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

	// Turn off character echo hide the cursor and disable input buffering
	gc.Echo(false)

	// Hide cursor
	gc.Cursor(0)

	// Disable input buffering
	gc.CBreak(true)

	// Set colors
	if err := gc.InitPair(1, gc.C_BLACK, gc.C_GREEN); err != nil {
		log.Fatal("InitPair failed: ", err)
	}

	stdscr.ColorOn(1)
	stdscr.Print("Use vim bindings to move the snake. Press 'q' to exit")
	stdscr.ColorOff(1)

	stdscr.Refresh()

	// Use maximum screen width
	rows, cols := stdscr.MaxYX()

	// Define object dimensions
	height, width := 2, 8
	y, x := rows/2, cols/2

	// Create a rectangle window that is a placeholder for the snake
	var snake *gc.Window
	snake, err = gc.NewWindow(height, width, y, x)
	if err != nil {
		log.Fatal(err)
	}

	// Init snake
	snake.MoveWindow(y, x)
	snake.Print(snake_ascii)
	snake.Refresh()

	// Init food
	spawnFood(stdscr)
	stdscr.Refresh()

	snake.Timeout(100)

	// Wait for keyboard input
	snake.Keypad(true)

	// Define timings
	c := time.NewTicker(time.Second / 2)
	c2 := time.NewTicker(time.Second / 4)

loop:
	for {
		stdscr.Refresh()

		select {
		case <-c.C:
			spawnFood(stdscr)
		case <-c2.C:
			spawnFood(stdscr)
		default:
			if !handleInput(stdscr, snake) {
				break loop
			}
		}

		// Flush characters that have changed
		gc.Update()

	}
	snake.Delete()
}
