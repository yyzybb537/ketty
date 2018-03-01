package ketty

import (
	"fmt"
	"sync"
	"sync/atomic"
	"golang.org/x/net/context"
	U "github.com/yyzybb537/ketty/url"
)

type RobinBalancer struct {
	addrs []U.Url
	chooseIndex uint32
	mtx sync.RWMutex
}

func init() {
	RegBalancer("", new(RobinBalancer))
	RegBalancer("default", new(RobinBalancer))
}

func (this *RobinBalancer) Filte(in []U.Url) (out []U.Url) {
	out = in
	return
}

func (this *RobinBalancer) Up(addr U.Url) (down func()) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.addrs = append(this.addrs, addr)
	return func() {
		this.mtx.Lock()
		defer this.mtx.Unlock()
		for i, v := range this.addrs {
			if v == addr {
				copy(this.addrs[i:], this.addrs[i+1:])
				this.addrs = this.addrs[:len(this.addrs)-1]
				return
            }
        }
    }
}

func (this *RobinBalancer) Get(ctx context.Context) (addr U.Url, put func(), err error) {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	if len(this.addrs) == 0 {
		return U.Url{}, nil, fmt.Errorf("No estab connection")
    }

	index := atomic.AddUint32(&this.chooseIndex, 1)
	index = index % uint32(len(this.addrs))
	return this.addrs[index], nil, nil
}

func (this *RobinBalancer) Clone() Balancer {
	return new(RobinBalancer)
}

