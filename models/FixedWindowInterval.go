package models

import "time"

// FixedWindowInterval :
type FixedWindowInterval struct {
	StartTime time.Time
	EndTime   time.Time
	Interval  time.Duration
}

func (fixedWindow *FixedWindowInterval) setStartAndEndTime() {
	currTime := time.Now().UTC()
	fixedWindow.StartTime = currTime
	fixedWindow.EndTime = currTime.Add(fixedWindow.Interval)
}

// Run :
func (fixedWindow *FixedWindowInterval) Run(cb func()) {
	go func() {
		ticker := time.NewTicker(fixedWindow.Interval)
		fixedWindow.setStartAndEndTime()
		for range ticker.C {
			cb()
			fixedWindow.setStartAndEndTime()
		}
	}()
}
