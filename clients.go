package ketty

import (
	"fmt"
	"sync"
	"time"
	"golang.org/x/net/context"
)

type UniqMap map[string]bool

// Implement Client interface, and manage multiply clients.
type Clients struct {
	AopList

	url       Url
	balancer  Balancer

	// manage all of address
	addrMtx   sync.Mutex
	addrs     map[string]Url

	// retry
	q		  chan interface{}

	// close
	closeMtx  sync.Mutex
	onClose   []func()
	closed    bool

	// root reference
	root	  *Clients
}

func newClients(url Url, balancer Balancer, root *Clients) *Clients {
	url.MetaData = nil
	c := &Clients{
		url : url,
		balancer : balancer,
		q : make(chan interface{}),
		closed : false,
		root : root,
    }
	if c.root == nil {
		c.root = c
    }
	c.onClose = append(c.onClose, func(){ close(c.q) })
	return c
}

func (this *Clients) dial() error {
	proto, err := GetProtocol(this.url.Protocol)
	if err == nil {
		return this.dialProtocol(proto)
    }

	driver, err := GetDriver(this.url.Protocol)
	if err == nil {
		return this.dialDriver(driver)
	}

	return fmt.Errorf("Error url, unkown protocol. url:%s", this.url.ToString())
}

func (this *Clients) dialProtocol(proto Protocol) error {
	client, err := proto.Dial(this.url)
	if err != nil {
		return err
	}

	this.url.MetaData = client
	this.balancer.Up(this.url)
	return nil
}

func (this *Clients) dialDriver(driver Driver) error {
	upC, downC, stop, err := driver.Watch(this.url)
	if err != nil {
		return err
	}
	this.onClose = append(this.onClose, stop)

	go func() {
		for {
			up := []Url{}
			down := []Url{}
			end := false
			select {
			case up, end = <-upC:
				break
			case down, end = <-downC:
				break
			}
			if end {
				break
			}

			// down
			for _, url := range down {
				key := url.ToString()
				this.addrMtx.Lock()
				var exists bool
				if url, exists = this.addrs[key]; !exists {
					this.addrMtx.Unlock()
					continue
                }

				// 清除
				delete(this.addrs, key)
				this.addrMtx.Unlock()

				// Close
				url.MetaData.(Client).Close()
            }

			// up
			up = this.balancer.Filte(up)
			for _, url := range up {
				key := url.ToString()
				this.addrMtx.Lock()
				var exists bool
				if url, exists = this.addrs[key]; exists {
					this.addrMtx.Unlock()
					continue
                }
				client := newClients(url, this.balancer, this.root)
				url.MetaData = client
				this.addrs[key] = url
				this.addrMtx.Unlock()

				// connect
				err = client.dial()
				if err != nil {
					down := this.balancer.Up(url)
					client.onClose = append(client.onClose, down)
					continue
                }
				
				// 连接失败, 转入后台重试
				go func() {
					err := client.retryDial()
					if err != nil {
						down := this.balancer.Up(url)
						client.onClose = append(client.onClose, down)
					}
				}()
			}
        }
	}()

	return nil
}

func (this *Clients) retryDial() error {
	for {
		select {
		case <-this.q:
			return fmt.Errorf("Stop retry dial")
		case <-time.After(time.Second * 3):
        }

		err := this.dial()
		if err != nil {
			continue
		}

		this.closeMtx.Lock()
		if this.closed {
			this.closeMtx.Unlock()
			this.Close()
			return fmt.Errorf("Client is closed")
		}
		this.closeMtx.Unlock()
		return nil
    }
}

func (this *Clients) Close() {
	this.closeMtx.Lock()
	onClose := this.onClose
	this.onClose = nil
	this.closed = true
	this.closeMtx.Unlock()
	for _, f := range onClose {
		f()
    }
}

func (this *Clients) Invoke(ctx context.Context, handle ServiceHandle, method string, req, rsp interface{}) error {
	url, put, err := this.balancer.Get(ctx)
	if err != nil {
		return err
	}
	if put != nil {
		defer put()
    }
	client := url.MetaData.(Client)
	return client.Invoke(SetAop(ctx, this.root.GetAop()), handle, method, req, rsp)
}

