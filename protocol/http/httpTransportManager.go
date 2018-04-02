package http_proto

import (
	"fmt"
	"net/http"
	"time"
	"sync"
	"net"
	"crypto/tls"
)

type rwTimeoutConn struct {
	*net.TCPConn
	rwTimeout time.Duration
}

func (this *rwTimeoutConn) Read(b []byte) (int, error) {
	err := this.TCPConn.SetDeadline(time.Now().Add(this.rwTimeout))
	if err != nil {
		return 0, err
	}
	return this.TCPConn.Read(b)
}

func (this *rwTimeoutConn) Write(b []byte) (int, error) {
	err := this.TCPConn.SetDeadline(time.Now().Add(this.rwTimeout))
	if err != nil {
		return 0, err
	}
	return this.TCPConn.Write(b)
}

var DefaultTransportManager *TransportManager = &TransportManager{tss: make(map[string]*http.Transport)}

type TransportManager struct {
	tss			map[string]*http.Transport
	rwlock		sync.RWMutex
}

func (this *TransportManager) GetTransport(ConnectTimeout, RWTimeout, ResponseHeaderTimeout time.Duration) (ts *http.Transport){
	var ok bool
	timeoutName := fmt.Sprintf("%s-%s-%s",int(ConnectTimeout.Seconds()),int(RWTimeout.Seconds()),int(ResponseHeaderTimeout.Seconds()))

	this.rwlock.RLock()
	if ts, ok = this.tss[timeoutName];!ok {
		this.rwlock.RUnlock()
		this.rwlock.Lock()
		ts = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		Dial := func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, ConnectTimeout)
			if err != nil {
				return nil, err
			}
			if RWTimeout > 0 {
				return &rwTimeoutConn{
					TCPConn:   conn.(*net.TCPConn),
					rwTimeout: RWTimeout,
				}, nil
			} else {
				return conn, nil
			}
		}
		ts.Dial = Dial
		ts.ResponseHeaderTimeout = ResponseHeaderTimeout
		ts.MaxIdleConnsPerHost = 2000
		this.tss[timeoutName] = ts
		this.rwlock.Unlock()
    }else {
		this.rwlock.RUnlock()
    }
	return
}
