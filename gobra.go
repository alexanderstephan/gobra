package main

import (
	gc "github.com/rthornton128/goncurses"
)

func InitSnake(stdscr *gc.Window) {
	screen_y , screen_x := stdscr.MaxYX()

	for i := 0; i < start_size ; i++ {
		snake.PushFront(Segment{y: screen_y/2, x: screen_x/2+i})
	}
}

func setSnakeDir(input *gc.Window, y, x int) bool {
	// Get input from a dedicated window, otherwise stdscr would be blocked
	k := input.GetChar()

	// Define input handlers with interrupt condition
	switch byte(k) {
	case 'w':
		if d != South {
			d = North
		}
	case 'a':
		if d != East {
			d = West
		}
	case 's':
		if d != North {
			d = South
		}
	case 'd':
		if d != West {
			d = East
		}
	case 'q':
		return false
	}
	return true
}

func MoveSnake() {
	snake.Remove(snake.Back())

	// Increment or decrement last position
	head_y := snake.Front().Value.(Segment).y
	head_x := snake.Front().Value.(Segment).x

	switch d {
	case North:
		head_y--
	case South:
		head_y++
	case West:
		head_x--
	case East:
		head_x++
	}

	snake.PushFront(Segment{head_y, head_x})
}

func GrowSnake(size int) {
	for i := 0; i < size; i++ {
		tail_y := snake.Back().Value.(Segment).y
		tail_x := snake.Back().Value.(Segment).x

		switch d {
		case North:
			tail_y++
		case South:
			tail_y--
		case West:
			tail_x++
		case East:
			tail_x--
		}
		snake.PushBack(Segment{y:tail_y, x:tail_x})
	}


}

func RenderSnake(stdscr *gc.Window) {
	currentSegment := snake.Front()

	for currentSegment != nil {
		stdscr.MoveAddChar(currentSegment.Value.(Segment).y, currentSegment.Value.(Segment).x, gc.Char(snake_body))
		currentSegment = currentSegment.Next()
	}
}
