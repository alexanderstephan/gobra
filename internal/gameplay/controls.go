package gameplay

import (
	"fmt"

	gc "github.com/rthornton128/goncurses"
)

// Controls
var (
	keyUp    byte
	keyLeft  byte
	keyDown  byte
	keyRight byte
)

// HandleKeys handles keyboard input for controlling the snake and performing other actions.
func HandleKeys(input *gc.Window, stdscr *gc.Window, myFood *Food) bool {
	// Get input from a dedicated window, otherwise stdscr would be blocked
	// Define input handlers with interrupt condition.
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

func initKeybindings(isVim bool) {
	// Remap to vim like bindings.
	if isVim {
		fmt.Println(" is Vim")
		keyLeft = 'h'
		keyDown = 'j'
		keyUp = 'k'
		keyRight = 'l'
	} else {
		fmt.Println(" is not Vim")
		keyLeft = 'a'
		keyUp = 'w'
		keyDown = 's'
		keyRight = 'd'
	}
}
