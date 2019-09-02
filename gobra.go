package main

import (
	gc "github.com/rthornton128/goncurses"
	"log"
	"math/rand"
	"time"
)

const food_char = 'X'

const snake_ascii = "#"

type Direction int

type Food struct {
	*gc.Window
	eaten bool
	y, x  int
}

type SnakeSegment struct {
	y, x    int
	segment string
	next    *SnakeSegment
}

type Snake struct {
	*gc.Window
	length int
	start  *SnakeSegment
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

func (s *Snake) insertSegments(newSegment *SnakeSegment) {
	if s.length == 0 {
		s.start = newSegment
	} else {
		currentSegment := s.start

		// Traverse the list until the next node is empty
		for currentSegment.next != nil {
			currentSegment = currentSegment.next
		}
		// Append new segment
		currentSegment.next = newSegment
	}
	// In both cases increment by one
	s.length++
}

func (s *Snake) createSnake() {
	//mySnake := &Snake{}
	size := 7
	for i := 0; i < size; i++ {
		node := SnakeSegment{}
		if s.length == 0 {
			node = SnakeSegment{segment: "@", y: 100, x: 100}
		} else {
			node = SnakeSegment{segment: "=", y: 100, x: 100 + i}
		}
		s.insertSegments(&node)
	}
}

func (s *Snake) renderSnake(snake *gc.Window) {
	list := s.start
	snake_string := ""
	for list != nil {
		snake_string += list.segment
		list = list.next
	}
	snake.Print(snake_string)
}

func spawnFood(stdscr *gc.Window) *Food {
	val1, val2 := stdscr.MaxYX()
	y := rand.Intn(val1)
	x := rand.Intn(val2)
	stdscr.MoveAddChar(y, x, food_char)
	return &Food{eaten: false, y: y, x: x}
}

func setSnakeDir(stdscr *gc.Window, snake *gc.Window) bool {
	rows, cols := stdscr.MaxYX()
	y, x := snake.YX()
	k := snake.GetChar()

	switch byte(k) {
	case 'q':
		return false
	case 'h':
		if x > 0 {
			if d != East {
				d = West
			}
		}
	case 'l':
		if x < cols {
			if d != West {
				d = East
			}
		}
	case 'k':
		if y > 1 {
			if d != South {
				d = North
			}
		}
	case 'j':
		if y < rows {
			if d != North {
				d = South
			}
		}
	}
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
	height, width := 1, 8
	y, x := rows/2, cols/2

	// Create a rectangle window that is a placeholder for the snake
	var snake *gc.Window
	snake, err = gc.NewWindow(height, width, y, x)
	if err != nil {
		log.Fatal(err)
	}

	// Init snake
	newSnake := &Snake{}
	newSnake.createSnake()
	snake.MoveWindow(y, x)
	newSnake.renderSnake(snake)
	snake.Refresh()

	// Init food
	spawnFood(stdscr)

	// Threshold for timeout
	snake.Timeout(100)

	// Wait for keyboard input
	snake.Keypad(true)

	// Define timings
	c := time.NewTicker(time.Second * 5)
	c2 := time.NewTicker(time.Second / 16)

loop:
	for {
		// Time events
		select {
		case <-c.C:
			spawnFood(stdscr)
		case <-c2.C:
			snake.Erase()
			snake.Refresh()
			newSnake.renderSnake(snake)
			snake.MoveWindow(y, x)
			snake.Refresh()

			switch d {
			case North:
				y--
			case South:
				y++
			case West:
				x--
			case East:
				x++
			}
		default:
			if !setSnakeDir(stdscr, snake) {
				break loop
			}
		}

		// Flush characters that have changed
		gc.Update()

	}
	snake.Delete()
}
