package main

import (
  gc "github.com/rthornton128/goncurses"
  "log"
)

func main() {
  // Initialize goncurses
  // End is required to preserve terminal after execution
  stdscr, err := gc.Init()
  if err != nil {
    log.Fatal(err)
  }
  defer gc.End()
  // Turn off character echo, hide the cursor and disable input buffering
  gc.Echo(false)
  gc.CBreak(true)
  gc.Cursor(0)

  stdscr.Print("Use vim bindings to move the snake. Press 'q' to exit")
  stdscr.NoutRefresh()

  rows, cols  :=  stdscr.MaxYX()
  height, width := 2, 8
  y, x := (rows-height)/2, (cols-width)/2
  // Create a rectangle window that is a placeholder for the snake
  var win *gc.Window
  win, err = gc.NewWindow(height, width, y, x)
  if err != nil {
    log.Fatal(err)
  }

  // Wait for keyboard input
  win.Keypad(true)

main:
  for {
    // Prevent output to terminal
    win.Erase()
    win.NoutRefresh()
    // Move the window and redraw it
    win.MoveWindow(y, x)
    win.Box(gc.ACS_VLINE, gc.ACS_HLINE)
    win.NoutRefresh()
    // Flush characters that have changed
    gc.Update()

    switch win.GetChar() {
    case 'q':
      break main
    case 'h':
      if x > 0 {
        x--
      }
    case 'l':
      if x < cols-width {
        x++
      }
    case 'k':
      if y > 1 {
        y--
      }
    case 'j':
      if y < rows-height {
        y++
      }
    }
  }
  win.Delete()
}
