package workpool

import (
	"sync"
	"errors"
	"sync/atomic"
)

type f func() error
type sig struct{}

var (
	ErrInvalidSize = errors.New("invalid size")
	ErrPoolClosed  = errors.New("pool is closed")
)

type Pool struct {
	capacity int32				// pool capacity
	workers []*Worker      		// worker slice
	runningCount int32    		// running worker number
	lock sync.Mutex      		// modify Pool lock
	release chan sig
	once sync.Once
}

func NewPool(size int) (*Pool, error) {
	if size < 0 {
		return nil, ErrInvalidSize
	}

	p := &Pool{
		capacity:   int32(size),
	}

	return p, nil
}

func (p *Pool) Submit(task f) error {
	if len(p.release) > 0 {
		return ErrPoolClosed
	}

	//w := p.getWorker()
	//w.sendTask(task)
	return nil
}

func (p *Pool) getWorker() *Worker {
	var w *Worker
	waiting := false

	p.lock.Lock()

	p.lock.Unlock()

	return nil
}

// Running returns the number of the currently running goroutines
func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.runningCount))
}

// Release Closed this pool
func (p *Pool) Release() error {
	p.once.Do(func() {
		p.release <- sig{}
		running := p.Running()
		for i := 0; i < running; i++ {
			//p.getWorker().stop()
		}
	})
	return nil
}