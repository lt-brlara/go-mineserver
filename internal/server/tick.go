package server

import (
	"time"

	"github.com/blara/go-mineserver/internal/log"
)

const (
	TICK_DURATION_GOAL time.Duration = 50 * time.Millisecond
	TICK_RATE_GOAL = 1.0 / TICK_DURATION_GOAL
)

func (s *Server) tick() {

	for {
		tickStart := time.Now()

		// do things "per-tick"

		tickRate := useRaminingTime(tickStart)
		log.Trace("Tick completed", 
			"tick_duration", log.Fmt("%dms", time.Since(tickStart).Milliseconds()),
			"tick_rate", log.Fmt("%.2f", tickRate),
		)
	}
}

func useRaminingTime(t time.Time) float64 {
		tickRate := 1.0 / time.Since(t).Seconds()

		if tickRate > TICK_RATE_GOAL.Seconds() {
			tickRate = 20
			time.Sleep((time.Millisecond * 50) - time.Since(t))
		}

		return tickRate
}
