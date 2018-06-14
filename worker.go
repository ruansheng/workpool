package workpool

import (
	"sync/atomic"
)

type Worker struct {
	pool *Pool
	task chan f
}

func (w *Worker) Run() {
	go func() {
		// block waiting worker end
		for f := range w.task {
			if f == nil {
				atomic.AddInt32(&w.pool.runningCount, -1)
				return
			}
			f()
			w.pool.putWorker(w)
		}
	}()
}

// stop this worker
func (w *Worker) Stop() {
	w.SendTask(nil)
}

// notice this worker
func (w *Worker) SendTask(data f) {
	w.task <- data
}