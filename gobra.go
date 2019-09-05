package main 

import (
	gc "github.com/rthornton128/goncurses"
)

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

func (s *Snake) MoveSnake(newDir Direction) {
	// Init starting position
	currentSegment := s.start

	for currentSegment.next != nil {
		currentSegment.y = currentSegment.next.y
		currentSegment.x = currentSegment.next.x
		currentSegment = currentSegment.next
	}

	// Move currentSegment to the last element in the list
	for currentSegment.next != nil {
		currentSegment = currentSegment.next
	}

	// Increment or decremt last position
	switch newDir {
	case North:
		currentSegment.y--
	case South:
		currentSegment.y++
	case West:
		currentSegment.x--
	case East:
		currentSegment.x++
	}
	// Append node with new position
	//newNode := &Segment{y: currentSegment.y, x: currentSegment.x}
	//s.InsertSegments(newNode)

	currentSegment = s.start

	for currentSegment.next != nil {
		currentSegment.y = currentSegment.next.y
		currentSegment.x = currentSegment.next.x
		currentSegment = currentSegment.next
	}

	s.length++
}

func (s *Snake) CutTail() {
	if s.length <= 3 {
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

func (s *Snake) CutFront() {
	currentSegment := s.start

	for currentSegment.next != nil {
		currentSegment.y = currentSegment.next.y
		currentSegment.x = currentSegment.next.x
		currentSegment = currentSegment.next
	}
	s.length--
}

func (s *Snake) RenderSnake(stdscr *gc.Window) {
	currentSegment := s.start
	for currentSegment != nil {
		stdscr.MoveAddChar(currentSegment.y, currentSegment.x, gc.Char(snake_body))
		currentSegment = currentSegment.next
	}
}

