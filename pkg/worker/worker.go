package worker

import (
	"context"
	"time"
)

type Task interface {
	Exec()
	GetInterval() time.Duration
}

func Start(w Task, ctx context.Context) {
	ticker := time.NewTicker(w.GetInterval())
	go func() {
		for {
			select {
			case <-ticker.C:
				w.Exec()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
