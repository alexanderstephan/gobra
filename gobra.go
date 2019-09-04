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

var d Direction = North

func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d]
}
func (s *Snake) InsertSegments(newSegment *Segment) {
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

func (s *Snake) CutTail() {
	if s.length == 0 {
		return
	}
	var previousSegment *Segment
	currentSegment := s.start

	for currentSegment.next != nil {
		previousSegment = currentSegment
		currentSegment = currentSegment.next
	}

	// currentSegment = currentSegment.next not possible since previousSegment wouldn't be unused
	previousSegment.next = currentSegment.next

	s.length--
}

func (s *Snake) InitSnake(stdscr *gc.Window) {
	snake_pos_y, snake_pos_x := stdscr.MaxYX()

	for i := 0; i < start_size; i++ {
		node := Segment{y: snake_pos_y / 2, x: snake_pos_x/2 - i}
		s.InsertSegments(&node)
	}
}

func (s *Snake) PushFront(newDir Direction) {
	newSegment := s.start
	nextSegment := s.start.next

	switch newDir {
	case North:
		newSegment = &Segment{y: (nextSegment.y - 1), x: nextSegment.x}
	case South:
		newSegment = &Segment{y: (nextSegment.y + 1), x: nextSegment.x}
	case West:
		newSegment = &Segment{y: nextSegment.y, x: (nextSegment.x - 1)}
	case East:
		newSegment = &Segment{y: nextSegment.y, x: (nextSegment.x + 1)}
	}

	s.start = newSegment

	s.length++
}

func (s *Snake) RenderSnake(stdscr *gc.Window) {
	snake_segment := s.start
	for snake_segment != nil {
		stdscr.MoveAddChar(snake_segment.y, snake_segment.x, gc.Char(snake_body))
		snake_segment = snake_segment.next
	}
}

func (f *Food) spawnFood(stdscr *gc.Window) Food {
	val1, val2 := stdscr.MaxYX()
	food_y := rand.Intn(val1)
	food_x := rand.Intn(val2)
	stdscr.MoveAddChar(food_y, food_x, food_char)
	random_food := Food{y: food_y, x: food_x}
	return random_food
}

func setSnakeDir(stdscr *gc.Window, snake *gc.Window, y, x int) bool {
	rows, cols := stdscr.MaxYX()
	k := snake.GetChar()

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
	y, x := rows/2, cols/2

	// Create a rectangle window that is a placeholder for the snake
	var snake *gc.Window
	snake, err = gc.NewWindow(0, 0, y, x)
	if err != nil {
		log.Fatal(err)
	}

	snake.Refresh()

	// Init snake
	newSnake := &Snake{}
	newSnake.InitSnake(stdscr)
	//node := Segment{symbol: snake_body, y: 40, x: 40}
	//newSnake.pushFront(&node)

	// Init food
	var food_y, food_x int
	stdscr.Refresh()

	// Threshold for timeout
	snake.Timeout(100)

	// Wait for keyboard input
	snake.Keypad(true)

loop:
	for {
		// Clear screen
		stdscr.Refresh()
		stdscr.Erase()

		// Show controls
		stdscr.ColorOn(1)
		stdscr.MovePrint(0, 0, "Use vim bindings to move the snake. Press 'q' to exit")
		stdscr.MovePrint(1, 0, newSnake.length)
		snake_head := newSnake.start

		if food_y == 0 && food_x == 0 {
			food_y = rand.Intn(rows)
			food_x = rand.Intn(cols)
		}
		stdscr.MoveAddChar(food_y, food_x, food_char)

		if snake_head.y == food_y && snake_head.x == food_x {
			food_y = rand.Intn(rows)
			food_x = rand.Intn(cols)
			stdscr.MoveAddChar(food_y, food_x, food_char)
		}

		stdscr.MovePrint(2, 0, food_y, food_x)
		stdscr.ColorOff(1)

		newSnake.CutTail()

		// Update snake position
		/*switch d {
		case North:
			newSnake.PushFront(North)
		case South:
			newSnake.PushFront(South)
		case West:
			newSnake.PushFront(West)
		case East:
			newSnake.PushFront(East)
		}
		*/
		// Render snake with altered position
		newSnake.RenderSnake(stdscr)

		// setSnakeDir returns false on exit -> interrupt loop
		if !setSnakeDir(stdscr, snake, y, x) {
			break loop
		}

		// Refresh changes in screen buffer
		stdscr.Refresh()

		// Flush characters that have changed
		gc.Update()

	}
	snake.Delete()
}
