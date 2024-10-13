package gameplay

import (
	"container/list"
	"gobra/internal/tools"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	gc "github.com/rthornton128/goncurses"
)

var (
	globalScore int
	startTime   time.Time
	newTime     time.Time

	screen Screen

	// Create a rectangle window that is a placeholder for the snake.
	input      *gc.Window
	gobraASCII = []string{
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
)

const (
	MaxRows int = 100
	MaxCols int = 100
)

type Segment struct {
	y, x int
}

type Food struct {
	y, x int
}

type Screen struct {
	rows, cols int
}

func drawBorder(stdscr *gc.Window) {
	stdscr.ColorOn(3)
	stdscr.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	stdscr.ColorOff(3)
}

func drawLogo(stdscr *gc.Window, rows, cols int) {
	stdscr.ColorOn(3)
	for i := 0; i < len(gobraASCII); i++ {
		stdscr.MovePrint(rows/2+i-5, cols/2-20, gobraASCII[i])
	}
	stdscr.ColorOff(3)
}

func initColors() {
	// Set up colors.
	tools.Check(gc.InitPair(1, gc.C_GREEN, gc.C_BLACK))
	tools.Check(gc.InitPair(2, gc.C_RED, gc.C_BLACK))
	tools.Check(gc.InitPair(3, gc.C_YELLOW, gc.C_BLACK))
	tools.Check(gc.InitPair(4, gc.C_MAGENTA, gc.C_BLACK))
	tools.Check(gc.InitPair(5, gc.C_BLACK, gc.C_BLACK))
	tools.Check(gc.InitPair(6, gc.C_WHITE, gc.C_BLACK))
}

// printFood prints the symbol for food at a given position.
func printFood(stdscr *gc.Window, newFood *Food, rows, cols int) {
	stdscr.ColorOn(2)
	stdscr.MoveAddChar(newFood.y, newFood.x, foodChar)
	stdscr.ColorOff(2)
}

// InitFood initializes the foods position.
func initFood(stdscr *gc.Window, myFood *Food, rows, cols int) {
	if myFood.y == 0 && myFood.x == 0 {
		for !testFoodCollision(stdscr, myFood, rows, cols) {
			myFood.y = rand.Intn(rows)
			myFood.x = rand.Intn(cols)
		}
		// Start the timer for the score calculation.
		startTime = time.Now()
	}
}

func gameOver(menu *gc.Window, rows, cols int) {
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

func handleCollisions(stdscr *gc.Window, myFood *Food, rows, cols int) bool {
	snakeFront := snake.Front().Value.(Segment)

	// Detect food collision.
	if snakeFront.y == myFood.y && snakeFront.x == myFood.x {
		for !testFoodCollision(stdscr, myFood, rows, cols) {
			myFood.y = rand.Intn(rows)
			myFood.x = rand.Intn(cols)
		}

		GrowSnake(growRate)

		// Calculate score.
		newTime = time.Now()

		r, err := os.ReadFile("/tmp/score")
		tools.Check(err)
		prevScore, err := strconv.Atoi(string(r))

		if err != nil {
			prevScore = 0
		}

		globalScore += (int(newTime.Sub(startTime) / 10000)) / scoreMulti

		if globalScore > prevScore {
			d := []byte(strconv.Itoa(globalScore))
			highscore = true
			err := os.WriteFile("/tmp/score", d, 0644)
			tools.Check(err)
		}

		// Reset timer for next food collection.
		startTime = time.Now()
		return false
	}

	// Check if head is element of the body.
	// First body element is the one after the head.
	bodyElement := snake.Front().Next()

	for bodyElement != nil {
		if (snakeFront.y == bodyElement.Value.(Segment).y) && (snakeFront.x == bodyElement.Value.(Segment).x) {
			snakeActive = false
			// Interrupt for-loop.
			break
		}
		// Move to the next element.
		bodyElement = bodyElement.Next()
	}
	return true
}

func boundaryCheck(nobounds bool, rows int, cols int) {
	snakeFront := snake.Front().Value.(Segment)

	// Detect boundaries.
	if !(nobounds) {
		if (snakeFront.y > rows-2) || (snakeFront.y < 1) || (snakeFront.x > cols-2) || (snakeFront.x < 1) {
			snakeActive = false
		}
		return
	}
	if snakeFront.y > rows-2 {
		// Hit bottom border.
		snake.Remove(snake.Back())
		snake.PushFront(Segment{1, snakeFront.x})
	} else if snakeFront.y < 1 {
		// Hit top border.
		snake.Remove(snake.Back())
		snake.PushFront(Segment{rows - 2, snakeFront.x})
	} else if snakeFront.x > cols-2 {
		// Hit right border.
		snake.Remove(snake.Back())
		snake.PushFront(Segment{snakeFront.y, 1})
	} else if snakeFront.x < 1 {
		// Hit left border.
		snake.Remove(snake.Back())
		snake.PushFront(Segment{snakeFront.y, cols - 2})
	}
}

func handleSnake(stdscr *gc.Window, rows int, cols int) {
	// Render snake with altered position
	// Move snake by one cell in the new direction
	if snakeActive {
		MoveSnake()
		stdscr.ColorOn(1)
		RenderSnake(stdscr)
		stdscr.ColorOff(1)
	} else if !snakeActive {
		stdscr.ColorOn(6)
		drawBorder(stdscr)
		RenderSnake(stdscr)
		gameOver(stdscr, rows, cols)
		stdscr.ColorOff(6)
	}
}

func testFoodCollision(stdscr *gc.Window, myFood *Food, rows, cols int) bool {
	return !(stdscr.MoveInChar(myFood.y, myFood.x) != ' ' || myFood.y == 0 || myFood.x == 0 || myFood.y == rows || myFood.x == cols)
}

func NewGame(stdscr *gc.Window, myFood *Food) {
	// Revive the snake.
	snakeActive = true

	// Reset direction.
	d = East

	// Empty list.
	snake = list.New()

	// Trigger initial food spawn.
	myFood.y = 0
	myFood.x = 0

	// Set up snake in original position.
	InitSnake(stdscr)

	// Reset score.
	globalScore = 0
}

func calcBoardSize(stdscr *gc.Window) {
	// Use maximum screen width.
	screen.rows, screen.cols = stdscr.MaxYX()
	if screen.rows > MaxRows && screen.cols > MaxCols {
		screen.rows = MaxRows
		screen.cols = MaxCols
	}

}

func SetupGameBoard() *gc.Window {
	// Setup stdscr buffer.
	stdscr, err := gc.Init()
	calcBoardSize(stdscr)
	if err != nil {
		log.Fatal(err)
	}

	// End is required to preserve terminal after execution.
	defer gc.End()

	// Has the terminal the capability to use color?
	if !gc.HasColors() {
		log.Fatal("Require a color capable terminal")
	}

	// Initalize the use of color.
	tools.Check(gc.StartColor())

	gc.Echo(false)
	gc.Cursor(0)    // Hide cursor
	gc.CBreak(true) // Disable input buffering

	// Define colors.
	initColors()

	input, err = gc.NewWindow(0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	input.Refresh()

	// Welcome screen with logo and controls.
	drawLogo(stdscr, screen.rows, screen.cols)
	return stdscr
}
