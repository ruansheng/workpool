package workpool

import (
	"sync"
	"errors"
	"sync/atomic"
	"math"
)

type f func() error
type sig struct{}

var (
	ErrInvalidSize = errors.New("invalid size")
	ErrPoolClosed  = errors.New("pool is closed")
)

type Pool struct {
	capacity int32				// pool capacity
	runningCount int32    		// running worker number
	workers []*Worker      		// worker slice
	lock sync.Mutex      		// modify Pool lock
	release chan sig			// notice pool close itself
	freeSignal chan sig			// receive worker free signal
	once sync.Once
}

func NewPool(size int) (*Pool, error) {
	if size <= 0 {
		return nil, ErrInvalidSize
	}

	p := &Pool{
		capacity:   int32(size),
		freeSignal: make(chan sig, math.MaxInt32),
		release:    make(chan sig, 1),
	}

	return p, nil
}

func (p *Pool) Submit(task f) error {
	if len(p.release) > 0 {
		return ErrPoolClosed
	}

	w := p.getWorker()
	w.SendTask(task)
	return nil
}

func (p *Pool) getWorker() *Worker {
	var w *Worker
	waiting := false

	p.lock.Lock()
	workers := p.workers
	n := len(workers) - 1
	if n < 0 {
		if p.runningCount >= p.capacity {
			waiting = true
		} else {
			p.runningCount++
		}
	} else {
		w = workers[n]
		workers[n] = nil
		p.workers = workers[:n]
	}
	p.lock.Unlock()

	if waiting {
		<-p.freeSignal   // when pool is full, block waiting
		for {
			p.lock.Lock()
			workers = p.workers
			i := len(workers) - 1
			if i < 0 {
				p.lock.Unlock()
				continue
			}
			w = workers[i]
			workers[i] = nil
			p.workers = workers[:i]
			p.lock.Unlock()
			break
		}
	} else if w == nil{
		// create worker
		w = &Worker{
			pool: p,
			task: make(chan f),
		}
		w.Run()
	}

	return w
}

// append idle worker to pool
func (p *Pool) putWorker(w *Worker) {
	p.lock.Lock()
	p.workers = append(p.workers, w)
	p.lock.Unlock()

	// receive blocked can get worker
	p.freeSignal <- sig{}
}

// ReSize change the capacity of this pool
func (p *Pool) ReSize(size int) {
	if size < p.Cap() {
		diff := p.Cap() - size
		for i := 0; i < diff; i++ {
			p.getWorker().Stop()
		}
	} else if size == p.Cap() {
		return
	}
	atomic.StoreInt32(&p.capacity, int32(size))
}

// Running returns the number of the currently running goroutines
func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.runningCount))
}

// Free returns the available goroutines to work
func (p *Pool) Free() int {
	return int(atomic.LoadInt32(&p.capacity) - atomic.LoadInt32(&p.runningCount))
}

// Cap returns the capacity of this pool
func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

// Release Closed this pool
func (p *Pool) Release() error {
	p.once.Do(func() {
		p.release <- sig{}
		running := p.Running()
		for i := 0; i < running; i++ {
			p.getWorker().Stop()
		}
	})
	return nil
}