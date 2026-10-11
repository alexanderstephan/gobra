package main

import (
	"flag"
	"gobra/internal/gameplay"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var cfg gameplay.Config

	flag.BoolVar(&cfg.Vim, "v", false, "Enable vim bindings")
	flag.BoolVar(&cfg.DebugInfo, "d", false, "Print debug info")
	flag.BoolVar(&cfg.NoBounds, "n", false, "Free boundaries")
	flag.BoolVar(&cfg.Sound, "s", false, "Enable sound")
	flag.IntVar((*int)(&cfg.Difficulty), "l", int(gameplay.Normal), "Difficulty level: 1=easy, 2=normal, 3=hard, 4=insane")

	flag.Parse()

	// Setup signal handler.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		gameplay.Run = false // TODO: Cancellable context
	}()
	gameplay.Start(&cfg)
}
