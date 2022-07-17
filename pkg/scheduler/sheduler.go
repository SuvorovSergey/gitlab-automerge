package scheduler

import "time"

func DoEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		go f()
	}
}
