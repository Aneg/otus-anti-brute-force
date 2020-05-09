package worker

import (
	"github.com/Aneg/otus-anti-brute-force/internal/services"
	"time"
)

func NewReloaderMasks(lists []services.IpGuard, errChan chan error) *ReloaderMasks {
	return &ReloaderMasks{
		lists:   lists,
		errChan: errChan,
	}
}

type ReloaderMasks struct {
	lists   []services.IpGuard
	errChan chan error
}

func (w *ReloaderMasks) Exec() {
	for i := range w.lists {
		if err := w.lists[i].Reload(); err != nil {
			w.errChan <- err
		}
	}
}

func (w *ReloaderMasks) GetInterval() time.Duration {
	return time.Second * 5
}

func (w *ReloaderMasks) IsStopped() bool {
	return false
}
