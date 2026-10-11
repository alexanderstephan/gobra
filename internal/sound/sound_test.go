package sound

import (
	"testing"
	"time"
)

func TestPlayWithoutContextDoesNotPanic(t *testing.T) {
	context = nil
	player := Play(FreqA, 50*time.Millisecond)
	if player != nil {
		t.Fatal("Play should return nil when audio context is unavailable")
	}
}
