package main

import (
	"container/list"
	"flag"
	gc "github.com/rthornton128/goncurses"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// Snake init data
const food_char = 'X'
const snake_alive = 'O'
const snake_dead = '+'
const start_size = 5
const grow_rate = 3
const score_multi = 20

// Controls
var key_up = 'w'
var key_left = 'a'
var key_down = 's'
var key_right = 'd'

var snake = list.New()

var snake_active bool
var debug_info bool
var game_started bool
var vim bool
var nobounds bool = false

var global_score int
var start_time time.Time
var new_time time.Time

var gobra_ascii = []string {
	`   ___     ___   | |__    __ __   ____`,
	` /  _  |  / _ \  | '_ \  | '__/  / _  |`,
	`|  (_| |   (_)   | |_)   | |    | (_| |`,
	` \__,  |  \___/  |_.__/  |_|     \__,_|`,
	` |___ /`,
}

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

func SetDir(input *gc.Window, stdscr *gc.Window, rows, cols int, newFood *Food) bool {
	// Get input from a dedicated window, otherwise stdscr would be blocked
	// Define input handlers with interrupt condition
	switch input.GetChar() {
	case gc.Key(key_up):
		if d != South {
			d = North
		}
	case gc.Key(key_left):
		if d != East {
			d = West
		}
	case gc.Key(key_down):
		if d != North {
			d = South
		}
	case gc.Key(key_right):
		if d != West {
			d = East
		}
	case ' ':
		if snake_active == false {
			NewGame(stdscr, rows, cols, newFood)
		}
	case 'Q':
		return false
	}
	return true
}

func NewGame(stdscr *gc.Window, rows, cols int, newFood *Food) {
	snake_active = true
	snake = list.New()
	newFood.y = 0
	newFood.x = 0
	InitSnake(stdscr)
	global_score = 0
}

func GameOver(menu *gc.Window, rows, cols int) {
	snake_active = false
	menu.ColorOn(2)
	menu.MovePrint((rows/2)-3, (cols/2)-4, "GAME OVER")
	menu.ColorOff(2)
	menu.ColorOn(1)
	menu.MovePrint(rows/2+1, cols/2-12, "Press SPACE to play again")
	menu.ColorOff(1)

	menu.ColorOn(3)
	menu.MovePrint(rows/2+5, cols/2-(len(strconv.Itoa(global_score))/2), global_score)
	menu.ColorOff(3)
}

func BoundaryCheck(nobounds *bool, rows int, cols int) {
	// Detect boundaries
	if !(*nobounds) {
		if (snake.Front().Value.(Segment).y > rows-2) || (snake.Front().Value.(Segment).y < 1) || (snake.Front().Value.(Segment).x > cols-2) || (snake.Front().Value.(Segment).x < 0) {
			snake_active = false
		}
	} else {
		// Hit bottom border
		if snake.Front().Value.(Segment).y > rows-2 {
			snake.Remove(snake.Back())
			snake.PushFront(Segment{1, snake.Front().Value.(Segment).x})
		}
		// Hit top border
		if snake.Front().Value.(Segment).y < 1 {
			snake.Remove(snake.Back())
			snake.PushFront(Segment{rows - 2, snake.Front().Value.(Segment).x})
		}
		// Hit right border
		if snake.Front().Value.(Segment).x > cols-2 {
			snake.Remove(snake.Back())
			snake.PushFront(Segment{snake.Front().Value.(Segment).y, 1})
		}
		// Hit left border
		if snake.Front().Value.(Segment).x < 0 {
			snake.Remove(snake.Back())
			snake.PushFront(Segment{snake.Front().Value.(Segment).y, cols - 2})
		}
	}
}

func TestFood(stdscr *gc.Window, rows int, cols int, newFood *Food) bool {
	// Draw food
	if stdscr.MoveInChar(newFood.y, newFood.x) == ' ' {
		return true
	} else {
		return false
	}
}

func HandleSnake(stdscr *gc.Window, rows int, cols int) {
	// Render snake with altered position
	if snake_active == true {
		// Move snake by one cell in the new direction
		MoveSnake()
		stdscr.ColorOn(1)
		RenderSnake(stdscr)
		stdscr.ColorOff(1)
	}
	if snake_active == false {
		stdscr.ColorOn(6)
		RenderSnake(stdscr)
		GameOver(stdscr, rows, cols)
		stdscr.ColorOff(6)
	}
}

func main() {
	vim := flag.Bool("V", false, "Enable vim bindings")
	debug_info := flag.Bool("D", false, "Print debug info")
	nobounds := flag.Bool("N", true, "Free boundaries")
	flag.Parse()

	if *vim {
		key_up = 'k'
		key_left = 'h'
		key_down = 'j'
		key_right = 'l'
	}

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

	if err := gc.InitPair(3, gc.C_YELLOW, gc.C_BLACK); err != nil {
		log.Fatal("InitPair failed: ", err)
	}

	if err := gc.InitPair(4, gc.C_RED, gc.C_BLACK); err != nil {
		log.Fatal("InitPair failed: ", err)
	}

	if err := gc.InitPair(5, gc.C_BLACK, gc.C_BLACK); err != nil {
		log.Fatal("InitPair failed: ", err)
	}

	if err := gc.InitPair(6, gc.C_WHITE, gc.C_BLACK); err !=  nil {
		log.Fatal( "InitPair failed: ", err)
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

	var menu *gc.Window
	menu, err = gc.NewWindow(40, 100, rows/2, cols/2)

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
	var i int
	for i = 0; i < len(gobra_ascii); i++ {
		stdscr.ColorOn(3)
		stdscr.MovePrint(rows/2+i-3, cols/2-20, gobra_ascii[i])
		stdscr.ColorOff(3)
	}
	stdscr.ColorOn(1)
	stdscr.MovePrint(rows/2+i+1, cols/2-25, "Control the snake with 'WASD'. Press Shift + Q to quit")
	stdscr.Refresh()
	stdscr.ColorOff(1)
	time.Sleep(1* time.Second)

loop:
	for {
		// Clear screen
		stdscr.Refresh()
		stdscr.Erase()

		stdscr.Box(gc.ACS_VLINE, gc.ACS_HLINE)

		if *debug_info == true {
			frame_counter++
			stdscr.MovePrint(1, 1, "DEBUG:")
			stdscr.MovePrint(2, 1, frame_counter)
			stdscr.MovePrint(3, 1, d)
			stdscr.MovePrint(4, 1, snake.Front().Value.(Segment).y)
			stdscr.MovePrint(4, 4, snake.Front().Value.(Segment).x)
			stdscr.MovePrint(5, 1, newFood.y)
			stdscr.MovePrint(5, 4, newFood.x)
			stdscr.MovePrint(6, 1, rows)
			stdscr.MovePrint(6, 4, cols)
			stdscr.MovePrint(7, 1, rune(stdscr.MoveInChar(1, 20)))
		}

		// setSnakeDir returns false on exit -> interrupt loop
		if !SetDir(input, stdscr, snake.Front().Value.(Segment).y, snake.Front().Value.(Segment).x, newFood) {
			break loop
		}

		HandleSnake(stdscr, rows, cols)

		// Init food position
		if newFood.y == 0 && newFood.x == 0 {
			newFood = &Food{y: rand.Intn(rows), x: rand.Intn(cols)}
			start_time = time.Now()
		}

		// Detect food collision
		if snake.Front().Value.(Segment).y == newFood.y && snake.Front().Value.(Segment).x == newFood.x {
			newFood = &Food{y: rand.Intn(rows), x: rand.Intn(cols)}
			if !TestFood(stdscr, rows, cols, newFood) {
				newFood = &Food{rows/2-5, cols/2+5}
			}
			GrowSnake(grow_rate)

			// Calculate score
			new_time = time.Now()
			global_score += (int(new_time.Sub(start_time)/10000))/score_multi
			start_time = time.Now()
		}

		// Check if head is element of the body
		e := snake.Front().Next()
		for e != nil {
			if (snake.Front().Value.(Segment).y == e.Value.(Segment).y) && (snake.Front().Value.(Segment).x == e.Value.(Segment).x) {
				snake_active = false
				break
			}
			e = e.Next()
		}


		stdscr.ColorOn(4)
		stdscr.MoveAddChar(newFood.y, newFood.x, food_char)
		stdscr.ColorOff(4)

		BoundaryCheck(nobounds, rows, cols)

		// Draw box around the screen
		stdscr.ColorOn(3)
		stdscr.Box(gc.ACS_VLINE, gc.ACS_HLINE)
		stdscr.ColorOff(3)

		// Refresh changes in screen buffer
		stdscr.Refresh()
		menu.Refresh()

		// Flush characters that have changed
		gc.Update()
	}
	input.Delete()
}
