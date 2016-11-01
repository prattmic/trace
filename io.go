package trace

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// context is the global trace package context.
type context struct {
	// enabled indicates tracing is enabled. It may be read atomically.
	enabled uint32

	// buf holds pending, unwritten, events.
	buf chan ViewerEvent

	// stop tells the reading goroutine to exit.
	stop chan struct{}

	// done indicates that the reading goroutine has exited.
	done chan struct{}

	// start is the beginning of time.
	start time.Time

	// mu protects the below.
	mu sync.Mutex

	// w is where the trace is written to.
	w io.Writer
}

var ctx context

func Start(w io.Writer) {
	if atomic.LoadUint32(&ctx.enabled) != 0 {
		panic("Tracing already enabled")
	}

	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.buf = make(chan ViewerEvent, 100)
	ctx.stop = make(chan struct{})
	ctx.done = make(chan struct{})
	ctx.w = w

	// TODO(prattmic): error checking
	preamble := `{"traceEvents": [`
	ctx.w.Write([]byte(preamble))

	go read()

	ctx.start = time.Now()
	atomic.StoreUint32(&ctx.enabled, 1)
}

func Stop() {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	atomic.StoreUint32(&ctx.enabled, 0)

	close(ctx.buf)
	<-ctx.done
	close(ctx.done)
}

func read() {
	enc := json.NewEncoder(ctx.w)

	first := true
	for {
		e, ok := <-ctx.buf
		if !ok {
			break
		}

		if first {
			first = false
		} else {
			ctx.w.Write([]byte(","))
		}

		err := enc.Encode(e)
		if err != nil {
			panic(fmt.Sprintf("Error encoding: %v", err))
		}
	}

	// Drain buf.
	for e := range ctx.buf {
		ctx.w.Write([]byte(","))

		err := enc.Encode(e)
		if err != nil {
			panic(fmt.Sprintf("Error encoding: %v", err))
		}
	}

	ending := "]}"
	ctx.w.Write([]byte(ending))
	ctx.done <- struct{}{}
}
