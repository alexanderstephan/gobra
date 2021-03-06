package main

import (
	"container/list"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	gc "github.com/alexanderstephan/goncurses"
	"github.com/hajimehoshi/oto"
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
	`                  888                     `,
	`                  888                     `,
	`                  888                     `,
	` .d88b.   .d88b.  88888b.  888d888 8888b. `,
	`d88P"88b d88""88b 888 "88b 888P"      "88b`,
	`888  888 888  888 888  888 888    .d888888`,
	`Y88b 888 Y88..88P 888 d88P 888    888  888`,
	` "Y88888  "Y88P"  88888P"  888    "Y888888`,
	`     888                                  `,
	`Y8b d88P                                  `,
	`"Y88P"                                    `,
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
	menu.MovePrint(rows/2-1, cols/2-15, "You died. Better luck next time!")
	menu.ColorOff(2)

	menu.ColorOn(1)
	menu.MovePrint(rows/2+1, cols/2-13, "Press 'SPACE' to play again!")
	menu.ColorOff(1)

	menu.ColorOn(3)
	if highscore {
		menu.MovePrint(3, cols/2-6, "New Highscore")
	}
	menu.ColorOff(3)
}

func BoundaryCheck(nobounds *bool, rows int, cols int) {
	snakeFront := snake.Front().Value.(Segment)

	// Detect boundaries
	if !(*nobounds) {
		if (snakeFront.y > rows-2) || (snakeFront.y < 1) || (snakeFront.x > cols-2) || (snakeFront.x < 1) {
			snakeActive = false
		}
	} else {
		if snakeFront.y > rows-2 {
			// Hit bottom border
			snake.Remove(snake.Back())
			snake.PushFront(Segment{1, snakeFront.x})
		} else if snakeFront.y < 1 {
			// Hit top border
			snake.Remove(snake.Back())
			snake.PushFront(Segment{rows - 2, snakeFront.x})
		} else if snakeFront.x > cols-2 {
			// Hit right border
			snake.Remove(snake.Back())
			snake.PushFront(Segment{snakeFront.y, 1})
		} else if snakeFront.x < 1 {
			// Hit left border
			snake.Remove(snake.Back())
			snake.PushFront(Segment{snakeFront.y, cols - 2})
		}
	}
}

func HandleSnake(stdscr *gc.Window, rows int, cols int) {
	// Render snake with altered position
	// Move snake by one cell in the new direction
	if snakeActive {
		MoveSnake()
		stdscr.ColorOn(1)
		RenderSnake(stdscr)
		stdscr.ColorOff(1)
	} else if !snakeActive {
		stdscr.ColorOn(6)
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
	snakeFront := snake.Front().Value.(Segment)

	// Detect food collision
	if snakeFront.y == myFood.y && snakeFront.x == myFood.x {
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
	// First body element is the one after the head
	bodyElement := snake.Front().Next()

	for bodyElement != nil {
		if (snakeFront.y == bodyElement.Value.(Segment).y) && (snakeFront.x == bodyElement.Value.(Segment).x) {
			snakeActive = false
			// Interrupt for-loop
			break
		}
		// Move to the next element
		bodyElement = bodyElement.Next()
	}
	return true
}

func DrawBorder(stdscr *gc.Window) {
	stdscr.ColorOn(3)
	stdscr.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	stdscr.ColorOff(3)
}

func DrawLogo(stdscr *gc.Window, rows, cols int) {
	stdscr.ColorOn(3)
	for i := 0; i < len(gobraAscii); i++ {
		stdscr.MovePrint(rows/2+i-5, cols/2-20, gobraAscii[i])
	}
	stdscr.ColorOff(3)
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
	stdscr.Refresh()
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
