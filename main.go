package main

import (
	"container/list"
	"flag"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	gc "github.com/alexanderstephan/goncurses"
)

// Parameters to tweak the playing experience
const foodChar = 'X'
const snakeAlive = 'O'
const snakeDead = '+'
const snakeHead = '0'
const startSize = 5
const growRate = 3
const scoreMulti = 20

// Is the game running?
var run = true

// Objects
var snake = list.New()
var newFood = Food{}

// States
var snakeActive bool
var highscore bool

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	vim := flag.Bool("v", false, "Enable vim bindings")
	debugInfo := flag.Bool("d", false, "Print debug info")
	noBounds := flag.Bool("n", false, "Free boundaries")
	sound := flag.Bool("s", false, "Enable sound")

	flag.Parse()

	// Randomize pseudo random functions
	rand.Seed(time.Now().Unix())
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		run = false
	}()

	if _, err := os.Stat("/tmp/score"); os.IsNotExist(err) {
		d := []byte("0")
		err = ioutil.WriteFile("/tmp/score", d, 0644)
		check(err)
	}

	stdscr := SetupGameBoard()
	InitSnake(stdscr) // Create initial snake
	initControls(*vim)
	input.Timeout(100) // Threshold for timeout
	input.Keypad(true) // Wait for keyboard input
	snakeActive = true // Snake starts alive
	frameCounter := 0  // Init frame count

	stdscr.Refresh()
	time.Sleep(1 * time.Second)

	var scoreLength int
	for run {
		// Clear screen
		stdscr.Refresh()
		stdscr.Erase()

		// Draw box around the screen (for collision detection)
		drawBorder(stdscr)

		// Print debug Infos
		if *debugInfo {
			frameCounter++
			stdscr.MovePrint(1, 1, "DEBUG:")
			stdscr.MovePrint(2, 1, frameCounter)
			stdscr.MovePrint(3, 1, d)
			stdscr.MovePrint(4, 1, snake.Front().Value.(Segment).y)
			stdscr.MovePrint(4, 4, snake.Front().Value.(Segment).x)
			stdscr.MovePrint(5, 1, newFood.y)
			stdscr.MovePrint(5, 4, newFood.x)
			stdscr.MovePrint(6, 1, screen.rows)
			stdscr.MovePrint(6, 4, screen.cols)
			stdscr.MovePrint(7, 1, rune(stdscr.MoveInChar(0, 0)))
		}

		// setSnakeDir returns false on exit -> interrupt loop
		if !setDir(input, stdscr, &newFood) {
			break
		}

		// Determine food position if not set yet
		initFood(stdscr, &newFood, screen.rows, screen.cols)

		//stdscr.Refresh()

		// Display snake (alive or dead)
		handleSnake(stdscr, screen.rows, screen.cols)

		// Handle collisions
		if !handleCollisions(stdscr, &newFood, screen.rows, screen.cols) && *sound {
			play(freqA, 250*time.Millisecond)
		}

		// Check if snake hit boundaries, if desired ports the snake to the other side of the screen
		boundaryCheck(noBounds, screen.rows, screen.cols)

		// Render food symbol
		printFood(stdscr, &newFood, screen.rows, screen.cols)

		// Overwrite border once again
		drawBorder(stdscr)
		scoreLength = len(strconv.Itoa(globalScore))

		// Write score to border
		stdscr.ColorOn(4)
		stdscr.MovePrint(0, (screen.cols/2)-(scoreLength/2), globalScore)
		stdscr.ColorOff(4)
		stdscr.ColorOn(3)
		stdscr.MoveAddChar(0, (screen.cols/2)-(scoreLength/2)-1, '|')

		if scoreLength%2 == 0 {
			stdscr.MoveAddChar(0, (screen.cols/2)+(scoreLength/2), '|')
		} else {
			stdscr.MoveAddChar(0, (screen.cols/2)+(scoreLength/2)+1, '|')
		}

		stdscr.ColorOff(3)
		// Refresh changes in screen buffer
		stdscr.Refresh()
		// Flush characters that have changed
		gc.Update()
	}
	input.Delete()
}
