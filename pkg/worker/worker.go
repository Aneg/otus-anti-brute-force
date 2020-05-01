package worker

import (
	"time"
)

type Task interface {
	Exec()
	GetInterval() time.Duration
	IsStopped() bool
}

func Start(w Task) {
	ticker := time.NewTicker(w.GetInterval())
	go func() {
		for {
			if w.IsStopped() {
				ticker.Stop()
				break
			}
			<-ticker.C
			w.Exec()
		}
	}()
}
