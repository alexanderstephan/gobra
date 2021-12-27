package main

import (
	gc "github.com/alexanderstephan/goncurses"
)

// Controls
var keyUp byte
var keyLeft byte
var keyDown byte
var keyRight byte

func setDir(input *gc.Window, stdscr *gc.Window, myFood *Food) bool {
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
		if !snakeActive {
			NewGame(stdscr, myFood)
		}
	case 'q':
		return false
	}
	return true
}

func initControls(isVim bool) {
	// Remap to vim like bindings
	if isVim {
		keyLeft = 'h'
		keyDown = 'j'
		keyUp = 'k'
		keyRight = 'l'
	} else {
		keyLeft = 'a'
		keyUp = 'w'
		keyDown = 's'
		keyRight = 'd'
	}
}
