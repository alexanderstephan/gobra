package gameplay

import (
	"testing"
	"time"
)

func TestTickDelayForDifficulty(t *testing.T) {
	tests := []struct {
		name     string
		diff     Difficulty
		want     time.Duration
	}{
		{name: "easy", diff: Easy, want: 180 * time.Millisecond},
		{name: "normal", diff: Normal, want: 120 * time.Millisecond},
		{name: "hard", diff: Hard, want: 90 * time.Millisecond},
		{name: "insane", diff: Insane, want: 70 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TickDelayForDifficulty(tt.diff); got != tt.want {
				t.Fatalf("TickDelayForDifficulty(%v) = %v, want %v", tt.diff, got, tt.want)
			}
		})
	}
}

func TestPowerUpReward(t *testing.T) {
	if got := ScoreBonusForPowerUp(ScoreBoost); got != 25 {
		t.Fatalf("ScoreBonusForPowerUp(ScoreBoost) = %d, want 25", got)
	}
	if got := ScoreBonusForPowerUp(Shield); got != 0 {
		t.Fatalf("ScoreBonusForPowerUp(Shield) = %d, want 0", got)
	}
}

func TestPauseToggle(t *testing.T) {
	paused = false
	TogglePause()
	if !paused {
		t.Fatal("TogglePause() should start a pause")
	}
	TogglePause()
	if paused {
		t.Fatal("TogglePause() should resume when already paused")
	}
}
