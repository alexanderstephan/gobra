package main

import (
	"flag"
	"fmt"
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

	flag.Parse()

	// Setup signal handler.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		gameplay.Run = false // TODO: Cancel context
	}()
	fmt.Println(cfg.Vim)
	gameplay.Start(&cfg)
}
