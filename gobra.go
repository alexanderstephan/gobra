package main

import (
    gc "github.com/alexanderstephan/goncurses"
)

// Draw snake in the screen center pointing East
func InitSnake(stdscr *gc.Window) {
    screen_y, screen_x := stdscr.MaxYX()
    for i := 0; i < startSize; i++ {
	snake.PushFront(Segment{y: screen_y / 2, x: screen_x / 2 + i})
    }
}

func MoveSnake() {
    // Delete last element of the snake
    snake.Remove(snake.Back())

    // Read coordinates of the first snake segment
    head_y := snake.Front().Value.(Segment).y
    head_x := snake.Front().Value.(Segment).x

    // Increment or decrement last position according to direction
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

    // Insert head with new position
    snake.PushFront(Segment{head_y, head_x})
}

func GrowSnake(size int) {
    for i := 0; i < size; i++ {
	tailY := snake.Back().Value.(Segment).y
	tailX := snake.Back().Value.(Segment).x

	// Move segment in the opposite direction
	switch d {
	case North:
	    tailY++
	case South:
	    tailY--
	case West:
	    tailX++
	case East:
	    tailX--
	}

	// Insert segment at back with new position
	snake.PushBack(Segment{y: tailY, x: tailX})
    }
}

func RenderSnake(stdscr *gc.Window) {
    // Traverse list and draw every segment to the screen depending on the snake state
    currentSegment := snake.Front()
    for currentSegment != nil {
	if snakeActive {
	    stdscr.MoveAddChar(currentSegment.Value.(Segment).y, currentSegment.Value.(Segment).x, gc.Char(snakeAlive))
	} else {
	    stdscr.MoveAddChar(currentSegment.Value.(Segment).y, currentSegment.Value.(Segment).x, gc.Char(snakeDead))
	}
	currentSegment = currentSegment.Next()
    }

    // Attach head
    if snakeActive {
	stdscr.MoveAddChar(snake.Front().Value.(Segment).y, snake.Front().Value.(Segment).x, snakeHead)
    }
}
