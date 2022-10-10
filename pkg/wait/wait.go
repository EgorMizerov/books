package wait

import "time"

func WaitUntil(task func() bool, deadline time.Time, interval time.Duration) bool {
	for !task() && time.Now().Before(deadline) {
		time.Sleep(interval)
	}
	return task()
}
