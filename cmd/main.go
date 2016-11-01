package main

import (
	"log"
	"os"
	"time"

	"github.com/prattmic/trace"
)

func main() {
	f, err := os.Create("/tmp/trace/trace.out")
	if err != nil {
		log.Printf("Error creating trace: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	trace.Start(f)
	defer trace.Stop()

	trace.InstantEvent()
	e := trace.NewEvent()
	time.Sleep(1*time.Millisecond)
	trace.InstantEvent()
	e2 := e.SubEvent()
	time.Sleep(1*time.Millisecond)
	trace.InstantEvent()
	e2.End()
	e.End()
}
