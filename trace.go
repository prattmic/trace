package trace

import (
	"sync/atomic"
	"time"
)

func InstantEvent() {
	if atomic.LoadUint32(&ctx.enabled) == 0 {
		return
	}

	ctx.buf <- ViewerEvent{
		Name: "instant event",
		Phase: "n",
		Categories: "instant",
		ID: 1,
		Time: float64(time.Since(ctx.start).Nanoseconds() / 1000),
	}
}

type Event struct {
	id uint64
}

func NewEvent() Event {
	ctx.buf <- ViewerEvent{
		Name: "event",
		Phase: "b",
		Categories: "duration",
		ID: 2,
		Time: float64(time.Since(ctx.start).Nanoseconds() / 1000),
	}
	return Event{id: 2}
}

func (e Event) End() {
	ctx.buf <- ViewerEvent{
		Name: "event",
		Phase: "e",
		Categories: "duration",
		ID: e.id,
		Time: float64(time.Since(ctx.start).Nanoseconds() / 1000),
	}
}

func (e Event) SubEvent() Event {
	ctx.buf <- ViewerEvent{
		Name: "event",
		Phase: "b",
		Categories: "duration",
		ID: e.id,
		Time: float64(time.Since(ctx.start).Nanoseconds() / 1000),
	}
	return e
}
