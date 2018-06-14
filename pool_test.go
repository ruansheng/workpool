package workpool

import (
	"testing"
	"fmt"
	"runtime"
)

func echo() error {
	fmt.Println("hello world!")
	//time.Sleep(time.Second * 3)
	return nil
}

func TestPool(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	pool, err := NewPool(3)
	if err != nil {
		t.Error(err)
	}

	for i:=0; i < 10 ;i++  {
		pool.Submit(echo)
	}

	//time.Sleep(time.Second * 60)
}

func BenchmarkPool(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	pool, err := NewPool(3)
	if err != nil {
		b.Error(err)
	}

	for i:=0; i < b.N ;i++  {
		pool.Submit(echo)
	}

	//time.Sleep(time.Second * 60)
}