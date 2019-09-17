package main

import (
	"container/list"
	"flag"
	"fmt"
	"github.com/hajimehoshi/oto"
	gc "github.com/rthornton128/goncurses"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// Snake init data
const foodChar = 'X'
const snakeAlive = 'O'
const snakeDead = '+'
const snakeHead = '0'
const startSize = 5
const growRate = 3
const scoreMulti = 20

// Sound frequencies
const freqA = 300

// Controls
var keyUp = 'w'
var keyLeft = 'a'
var keyDown = 's'
var keyRight = 'd'

// Objects
var snake = list.New()
var newFood = Food{}

// States
var snakeActive bool
var highscore bool

// Trackers
var globalScore int
var startTime time.Time
var newTime time.Time

// Audio
var (
	sampleRate      = flag.Int("samplerate", 44100, "sample rate")
	channelNum      = flag.Int("channelnum", 3, "number of channels")
	bitDepthInBytes = flag.Int("bitdepthinbytes", 2, "bit depth in bytes")
)

// Logo
var gobraAscii = []string{
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

func SetDir(input *gc.Window, stdscr *gc.Window, myFood *Food) bool {
	// Get input from a dedicated window, otherwise stdscr would be blocked
	// Define input handlers with interrupt condition
	switch input.GetChar() {
	case gc.Key(keyUp):
		if d != South {
			d = North
		}
	case gc.Key(keyLeft):
		if d != East {
			d = West
		}
	case gc.Key(keyDown):
		if d != North {
			d = South
		}
	case gc.Key(keyRight):
		if d != West {
			d = East
		}
	case ' ':
		if snakeActive == false {
			NewGame(stdscr, myFood)
		}
	case 'Q':
		return false
	}
	return true
}

func NewGame(stdscr *gc.Window, myFood *Food) {
	// Revive the snake
	snakeActive = true

	// Reset direction
	d = East

	// Empty list
	snake = list.New()

	// Trigger initial food spawn
	myFood.y = 0
	myFood.x = 0

	// Set up snake in original position
	InitSnake(stdscr)

	// Reset score
	globalScore = 0
}

func GameOver(menu *gc.Window, rows, cols int) {
	snakeActive = false
	menu.ColorOn(2)
	menu.MovePrint((rows / 2), (cols/2)-4, "GAME OVER")
	menu.ColorOff(2)
	menu.ColorOn(1)
	menu.MovePrint(rows/2+3, cols/2-12, "Press 'SPACE' to play again")
	menu.ColorOff(1)

	menu.ColorOn(3)
	if highscore {
		menu.MovePrint(rows/2+4, cols/2-6, "New Highscore")
	}
	menu.ColorOff(3)
}

func BoundaryCheck(nobounds *bool, rows int, cols int) {
	// Detect boundaries
	if !(*nobounds) {
		if (snake.Front().Value.(Segment).y > rows-2) || (snake.Front().Value.(Segment).y < 1) || (snake.Front().Value.(Segment).x > cols-2) || (snake.Front().Value.(Segment).x < 1) {
			snakeActive = false
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
		if snake.Front().Value.(Segment).x < 1 {
			snake.Remove(snake.Back())
			snake.PushFront(Segment{snake.Front().Value.(Segment).y, cols - 2})
		}
	}
}

func HandleSnake(stdscr *gc.Window, rows int, cols int) {
	// Render snake with altered position
	if snakeActive == true {
		// Move snake by one cell in the new direction
		MoveSnake()
		stdscr.ColorOn(1)
		RenderSnake(stdscr)
		stdscr.ColorOff(1)
	}
	if snakeActive == false {
		stdscr.ColorOn(6)
		stdscr.Erase()
		stdscr.Refresh()
		DrawBorder(stdscr)
		RenderSnake(stdscr)
		GameOver(stdscr, rows, cols)
		stdscr.ColorOff(6)
	}
}

func PrintFood(stdscr *gc.Window, newFood *Food, rows, cols int) {
	stdscr.ColorOn(2)
	stdscr.MoveAddChar(newFood.y, newFood.x, foodChar)
	stdscr.ColorOff(2)
}

// Init food position
func InitFood(stdscr *gc.Window, myFood *Food, rows, cols int) {
	if myFood.y == 0 && myFood.x == 0 {
		for !TestFoodCollision(stdscr, myFood, rows, cols) {
			myFood.y = rand.Intn(rows)
			myFood.x = rand.Intn(cols)
		}
		// Start timer for score
		startTime = time.Now()
	}
}

func TestFoodCollision(stdscr *gc.Window, myFood *Food, rows, cols int) bool {
	if stdscr.MoveInChar(myFood.y, myFood.x) != ' ' || myFood.y == 0 || myFood.x == 0 || myFood.y == rows || myFood.x == cols {
		return false
	}
	return true
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func HandleCollisions(stdscr *gc.Window, myFood *Food, rows, cols int) bool {
	// Detect food collision
	if snake.Front().Value.(Segment).y == myFood.y && snake.Front().Value.(Segment).x == myFood.x {
		for !TestFoodCollision(stdscr, myFood, rows, cols) {
			myFood.y = rand.Intn(rows)
			myFood.x = rand.Intn(cols)
		}

		GrowSnake(growRate)

		// Calculate score
		newTime = time.Now()

		r, err := ioutil.ReadFile("/tmp/score")
		check(err)
		prevScore, err := strconv.Atoi(string(r))
		globalScore += (int(newTime.Sub(startTime) / 10000)) / scoreMulti

		if globalScore > prevScore {
			d := []byte(strconv.Itoa(globalScore))
			highscore = true
			err := ioutil.WriteFile("/tmp/score", d, 0644)
			check(err)
		}

		// Reset timer for next food collection
		startTime = time.Now()
		return false
	}

	// Check if head is element of the body
	e := snake.Front().Next()
	for e != nil {
		if (snake.Front().Value.(Segment).y == e.Value.(Segment).y) && (snake.Front().Value.(Segment).x == e.Value.(Segment).x) {
			snakeActive = false
			break
		}
		e = e.Next()
	}
	return true
}

func DrawBorder(stdscr *gc.Window) {
	stdscr.ColorOn(3)
	stdscr.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	stdscr.ColorOff(3)
}

func DrawLogo(stdscr *gc.Window, rows, cols int) {
	var i int
	for i = 0; i < len(gobraAscii); i++ {
		stdscr.ColorOn(3)
		stdscr.MovePrint(rows/2+i-3, cols/2-20, gobraAscii[i])
		stdscr.ColorOff(3)
	}
	stdscr.ColorOn(1)
	stdscr.MovePrint(rows/2+i+1, cols/2-25, "Control the snake with 'WASD'. Press Shift + Q to quit")
	stdscr.Refresh()
	stdscr.ColorOff(1)
}

func InitColors() {
	// Set up colors
	if err := gc.InitPair(1, gc.C_GREEN, gc.C_BLACK); err != nil {
		log.Fatal("InitPair failed: ", err)
	}

	if err := gc.InitPair(2, gc.C_RED, gc.C_BLACK); err != nil {
		log.Fatal("InitPair failed: ", err)
	}

	if err := gc.InitPair(3, gc.C_YELLOW, gc.C_BLACK); err != nil {
		log.Fatal("InitPair failed: ", err)
	}

	if err := gc.InitPair(4, gc.C_MAGENTA, gc.C_BLACK); err != nil {
		log.Fatal("InitPair failed: ", err)
	}

	if err := gc.InitPair(5, gc.C_BLACK, gc.C_BLACK); err != nil {
		log.Fatal("InitPair failed: ", err)
	}

	if err := gc.InitPair(6, gc.C_WHITE, gc.C_BLACK); err != nil {
		log.Fatal("InitPair failed: ", err)
	}
}

func main() {
	vim := flag.Bool("v", false, "Enable vim bindings")
	debugInfo := flag.Bool("d", false, "Print debug info")
	noBounds := flag.Bool("n", false, "Free boundaries")
	sound := flag.Bool("s", false, "Enable sound")
	flag.Parse()

	// Remap to vim like bindings
	if *vim {
		keyUp = 'k'
		keyLeft = 'h'
		keyDown = 'j'
		keyRight = 'l'
	}

	// Setup stdscr buffer
	stdscr, err := gc.Init()

	// Use maximum screen width
	rows, cols := stdscr.MaxYX()

	if rows < 15 || cols < 20 {
		fmt.Printf("Screen resolution too small")
		return
	}

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

	// Define colors
	InitColors()

	// Create a rectangle window that is a placeholder for the snake
	var input *gc.Window
	input, err = gc.NewWindow(0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	input.Refresh()

	if _, err := os.Stat("/tmp/score"); os.IsNotExist(err) {
		d := []byte("0")
		err = ioutil.WriteFile("/tmp/score", d, 0644)
		check(err)
	}

	// Init sounds
	c, err := oto.NewContext(*sampleRate, *channelNum, *bitDepthInBytes, 4016)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup

	// Create initial snake
	InitSnake(stdscr)

	// Snake starts alive
	snakeActive = true

	// Init frame count
	frameCounter := 0

	// Threshold for timeout
	input.Timeout(100)

	// Wait for keyboard input
	input.Keypad(true)

	// Welcome screen with logo and controls
	DrawLogo(stdscr, rows, cols)
	time.Sleep(1 * time.Second)

loop:
	for {
		// Clear screen
		stdscr.Refresh()
		stdscr.Erase()

		// Draw box around the screen (for collision detection)
		DrawBorder(stdscr)

		// Print debug Infos
		if *debugInfo == true {
			frameCounter++
			stdscr.MovePrint(1, 1, "DEBUG:")
			stdscr.MovePrint(2, 1, frameCounter)
			stdscr.MovePrint(3, 1, d)
			stdscr.MovePrint(4, 1, snake.Front().Value.(Segment).y)
			stdscr.MovePrint(4, 4, snake.Front().Value.(Segment).x)
			stdscr.MovePrint(5, 1, newFood.y)
			stdscr.MovePrint(5, 4, newFood.x)
			stdscr.MovePrint(6, 1, rows)
			stdscr.MovePrint(6, 4, cols)
			stdscr.MovePrint(7, 1, rune(stdscr.MoveInChar(0, 0)))
		}

		// setSnakeDir returns false on exit -> interrupt loop
		if !SetDir(input, stdscr, &newFood) {
			break loop
		}

		// Determine food position if not set yet
		InitFood(stdscr, &newFood, rows, cols)

		//stdscr.Refresh()

		// Display snake (alive or dead)
		HandleSnake(stdscr, rows, cols)

		// Handle collisions
		if !HandleCollisions(stdscr, &newFood, rows, cols) && *sound {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := Play(c, freqA, 100*time.Millisecond); err != nil {
					panic(err)
				}
			}()
		}

		// Check if snake hit boundaries, if desired ports the snake to the other side of the screen
		BoundaryCheck(noBounds, rows, cols)

		// Render food symbol
		PrintFood(stdscr, &newFood, rows, cols)

		// Overwrite border once again
		DrawBorder(stdscr)

		// Write score to border
		stdscr.ColorOn(4)
		stdscr.MovePrint(0, cols/2-(len(strconv.Itoa(globalScore))/2), globalScore)
		stdscr.ColorOff(4)
		stdscr.ColorOn(3)
		stdscr.MoveAddChar(0, cols/2-1-(len(strconv.Itoa(globalScore))/2), '|')
		if len(strconv.Itoa(globalScore))%2 == 0 {
			stdscr.MoveAddChar(0, cols/2+(len(strconv.Itoa(globalScore))/2), '|')
		} else {
			stdscr.MoveAddChar(0, cols/2+1+(len(strconv.Itoa(globalScore))/2), '|')
		}
		stdscr.ColorOff(3)

		// Refresh changes in screen buffer
		stdscr.Refresh()
		// Flush characters that have changed
		gc.Update()
	}
	wg.Wait()
	c.Close()
	input.Delete()
}
