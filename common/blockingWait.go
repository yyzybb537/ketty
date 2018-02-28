package common

import (
	"sync"
)

type BlockingWait struct {
	mtx		sync.Mutex
	queue	chan interface{}
	count	int64
	wakeup  bool
}

func NewBlockingWait() *BlockingWait {
	return &BlockingWait{
		queue : make(chan interface{}, 1024),
    }
}

func (this *BlockingWait) Wait() {
	this.mtx.Lock()
	if this.wakeup {
		this.mtx.Unlock()
		return
    }

	this.count++
	this.mtx.Unlock()
	<-this.queue
}

func (this *BlockingWait) Notify() {
	if this.wakeup {
		return
    }

	this.mtx.Lock()
	defer this.mtx.Unlock()
	if this.wakeup {
		return
    }

	this.wakeup = true
	var i int64
	for ; i < this.count; i++ {
		this.queue <- nil
    }
	this.count = 0
}

