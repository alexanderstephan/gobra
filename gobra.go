package main

import (
	gc "github.com/rthornton128/goncurses"
	"log"
	"math/rand"
	"time"
)

var snake_ascii = []string{
	`====@`,
}

var food_ascii = []string{
	`X`,
}

type Object interface {
	Cleanup()
	Collide(int)
	Draw(*gc.Window)
	Expired(int, int) bool
	Update()
}

type Food struct {
	*gc.Window
	eaten bool
}

type Snake struct {
	//*gc.Window
	alive     bool
	y, x      int
	segements int
}

var objects = make([]Object, 0, 16)

func spawnFood(stdscr *gc.Window) *Food {
	val1, val2 := stdscr.MaxYX()
	y := rand.Intn(val1)
	x := rand.Intn(val2)
	stdscr.MovePrint(x, y, food_ascii)
	return &Food{eaten: false}
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
	spawnFood(stdscr)
	stdscr.Refresh()

	// Use maximum screen width
	rows, cols := stdscr.MaxYX()

	// Define object dimensions
	height, width := 2, 8
	y, x := (rows-height)/2, (cols-width)/2

	// Create a rectangle window that is a placeholder for the snake
	var win *gc.Window
	win, err = gc.NewWindow(height, width, y, x)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for keyboard input
	win.Keypad(true)

main:
	for {
		// Prevent output to terminal
		win.Erase()
		//stdscr.Refresh()
		win.Refresh()

		// stdscr.GetChar()
		//stdscr.SetBackground(gc.Char('x')) //| gc.ColorPair(1))
		// stdscr.ColorOn(1)

		// Move the window and redraw it
		win.MoveWindow(y, x)
		win.ColorOn(1)
		win.Print(snake_ascii)
		win.ColorOff(1)
		//win.Box(gc.ACS_VLINE, gc.ACS_HLINE)

		win.Refresh()

		// Flush characters that have changed
		gc.Update()

		// Get input and manipulate object position
		switch win.GetChar() {
		case 'q':
			break main
		case 'h':
			if x > 0 {
				x--
			}
		case 'l':
			if x < cols-width {
				x++
			}
		case 'k':
			if y > 1 {
				y--
			}
		case 'j':
			if y < rows-height {
				y++
			}
		}
	}
	win.Delete()
}
