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
	segment rune
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
	size := 7
	for i := 0; i < size; i++ {
		node := SnakeSegment{}
		if s.length == 0 {
			node = SnakeSegment{segment: '@', y: 20, x: 20}
		} else {
			node = SnakeSegment{segment: '=', y: 20, x: 20 + i}
		}
		s.insertSegments(&node)
	}
}

func (s *Snake) renderSnake(stdscr *gc.Window) {
	list := s.start
	for list != nil {
		stdscr.MoveAddChar(list.y, list.x, gc.Char(list.segment))
		list = list.next
	}
}

func spawnFood(stdscr *gc.Window) *Food {
	val1, val2 := stdscr.MaxYX()
	y := rand.Intn(val1)
	x := rand.Intn(val2)
	stdscr.MoveAddChar(y, x, food_char)
	return &Food{eaten: false, y: y, x: x}
}

func (s *Snake) setSnakeDir(stdscr *gc.Window) bool {
	rows, cols := stdscr.MaxYX()
	snake_pos := s.start
	y, x := snake_pos.y, snake_pos.x
	k := stdscr.GetChar()

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

	// Use maximum screen width
	rows, cols := stdscr.MaxYX()

	// Define object dimensions
	// height, width := 1, 8
	y, x := rows/2, cols/2

	// Create a rectangle window that is a placeholder for the snake
	/*var snake *gc.Window
	snake, err = gc.NewWindow(height, width, y, x)
	if err != nil {
		log.Fatal(err)
	}
	*/
	stdscr.MovePrint(y, x-10, "Press any key to start")

	// Init snake
	newSnake := &Snake{}
	newSnake.createSnake()

	// Threshold for timeout
	// snake.Timeout(100)

	// Wait for keyboard input
	stdscr.Keypad(true)

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
			stdscr.Erase()
			stdscr.Refresh()
			stdscr.ColorOn(1)
			stdscr.MovePrint(0, 0, "Use vim bindings to move the snake. Press 'q' to exit")
			newSnake.renderSnake(stdscr)
			stdscr.ColorOff(1)
			stdscr.Refresh()

			switch d {
			case North:
				mySegment := newSnake.start
				for mySegment != nil {
					mySegment.y--
					mySegment = mySegment.next
				}
			case South:
				mySegment := newSnake.start
				for mySegment != nil {
					mySegment.y++
					mySegment = mySegment.next
				}
			case West:
				mySegment := newSnake.start
				for mySegment != nil {
					mySegment.x--
				}
			case East:
				mySegment := newSnake.start
				for mySegment != nil {
					mySegment.x++
				}
			}
		default:
			if !newSnake.setSnakeDir(stdscr) {
				break loop
			}
		}

		// Flush characters that have changed
		gc.Update()

	}
	stdscr.Delete()
}
