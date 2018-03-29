package ketty

import (
	"fmt"
	"sync"
	"time"
	"golang.org/x/net/context"
	"github.com/yyzybb537/ketty/common"
	U "github.com/yyzybb537/ketty/url"
	P "github.com/yyzybb537/ketty/protocol"
	O "github.com/yyzybb537/ketty/option"
	B "github.com/yyzybb537/ketty/balancer"
	D "github.com/yyzybb537/ketty/driver"
	A "github.com/yyzybb537/ketty/aop"
	COM "github.com/yyzybb537/ketty/common"
)

// Implement Client interface, and manage multiply clients.
type Clients struct {
	A.AopList

	url       U.Url
	balancer  B.Balancer
	opt       O.OptionI

	// manage all of address
	addrMtx   sync.Mutex
	addrs     map[string]U.Url

	// retry
	q		  chan interface{}

	// close
	closeMtx  sync.Mutex
	onClose   []func()
	closed    bool

	// blocking wait queue
	blockingWait *common.BlockingWait

	// root reference
	root	  *Clients
}

func newClients(url U.Url, balancer B.Balancer, root *Clients) *Clients {
	url.MetaData = nil
	c := &Clients{
		url : url,
		balancer : balancer,
		addrs : map[string]U.Url{},
		q : make(chan interface{}),
		closed : false,
		blockingWait : common.NewBlockingWait(),
		root : root,
    }
	if c.root == nil {
		c.root = c
    }
	c.onClose = append(c.onClose, func(){ close(c.q) })
	return c
}

func (this *Clients) dial() error {
	//GetLog().Debugf("dial(%s)", this.url.ToString())
	proto, err := P.GetProtocol(this.url.GetMainProtocol())
	if err == nil {
		return this.dialProtocol(proto)
    }

	driver, err := D.GetDriver(this.url.Protocol)
	if err == nil {
		return this.dialDriver(driver)
	}

	return fmt.Errorf("Error url, unkown protocol. url:%s", this.url.ToString())
}

func (this *Clients) dialProtocol(proto P.Protocol) error {
	client, err := proto.Dial(this.url)
	if err != nil {
		return err
	}
	if this.opt != nil {
		client.SetOption(this.opt)
	}

	this.url.MetaData = client
	this.balancer.Up(this.url)
	return nil
}

func (this *Clients) dialDriver(driver D.Driver) error {
	upC, downC, stop, err := driver.Watch(this.url)
	if err != nil {
		return err
	}
	this.onClose = append(this.onClose, stop)

	go func() {
		for {
			up := []U.Url{}
			down := []U.Url{}
			readOk := true
			select {
			case up, readOk = <-upC:
				break
			case down, readOk = <-downC:
				break
			}
			if !readOk {
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
				url.MetaData.(P.Client).Close()
            }

			// up
			up = this.balancer.Filte(up)
			for _, url := range up {
				key := url.ToString()
				this.addrMtx.Lock()
				var exists bool
				if _, exists = this.addrs[key]; exists {
					this.addrMtx.Unlock()
					continue
                }
				client := newClients(url, this.balancer, this.root)
				if this.opt != nil {
					client.SetOption(this.opt)
				}
				url.MetaData = client
				this.addrs[key] = url
				this.addrMtx.Unlock()

				// connect
				err = client.dial()
				if err == nil {
					down := this.balancer.Up(url)
					client.onClose = append(client.onClose, down)
					this.blockingWait.Notify()
					continue
                }
				
				// 连接失败, 转入后台重试
				go func() {
					err := client.retryDial()
					if err != nil {
						down := this.balancer.Up(url)
						client.onClose = append(client.onClose, down)
						this.blockingWait.Notify()
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

func (this *Clients) SetOption(opt O.OptionI) error {
	this.opt = opt
	return nil
}

func (this *Clients) Invoke(ctx context.Context, handle COM.ServiceHandle, method string, req, rsp interface{}) error {
	url, put, err := this.balancer.Get(ctx)
	if err != nil {
		this.blockingWait.Wait()
		url, put, err = this.balancer.Get(ctx)
		if err != nil {
			return err
		}
	}
	if put != nil {
		defer put()
    }
	client := url.MetaData.(P.Client)
	return client.Invoke(A.SetAop(ctx, this.root.GetAop()), handle, method, req, rsp)
}

