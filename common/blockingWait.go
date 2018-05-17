package common

import (
	"sync"
	"time"
	"github.com/pkg/errors"
)

type BlockingWait struct {
	mtx		sync.Mutex
	queue	chan interface{}
	count	int64
	waitIdx int64
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

func (this *BlockingWait) TimedWait(dur time.Duration) error {
	if dur == 0 {
		this.Wait()
		return nil
	}

	this.mtx.Lock()
	if this.wakeup {
		this.mtx.Unlock()
		return nil
    }

	idx := this.waitIdx
	this.count++
	this.mtx.Unlock()
	select {
	case <-this.queue:
		return nil
	case <-time.NewTimer(dur).C:
		this.mtx.Lock()
		defer this.mtx.Unlock()
		if idx == this.waitIdx && this.count > 0 {
			this.count--
		}
		return errors.Errorf("wait timeout")
	}
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
	this.waitIdx++
	var i int64
	for ; i < this.count; i++ {
		this.queue <- nil
    }
	this.count = 0
}

