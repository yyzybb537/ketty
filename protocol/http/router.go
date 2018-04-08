package http_proto

import (
	U "github.com/yyzybb537/ketty/url"
	"github.com/yyzybb537/ketty/log"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"sync"
	"time"
	"github.com/yyzybb537/gls"
)

type Router struct {
	http.Server
	mux *http.ServeMux
	addr string
	proto string
	served bool
	opt        *HttpOption
}

var gRouters = map[string]*Router{}
var gMtx sync.Mutex

func getRouter(addr string) *Router {
	gMtx.Lock()
	defer gMtx.Unlock()

	if r, exists := gRouters[addr]; exists {
		return r
    }

	r := newRouter(addr)
	gRouters[addr] = r
	return r
}

func newRouter(addr string) *Router {
	r := new(Router)
	r.mux = http.NewServeMux()
	r.Handler = r.mux
	r.addr = addr
	return r
}

func (this *Router) Register(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	this.mux.HandleFunc(pattern, handler)
}

func (this *Router) RServe(proto string) error {
	if this.served {
		return nil
    }

	if this.opt != nil {
		this.Server.WriteTimeout = time.Duration(this.opt.Option.TimeoutMilliseconds) * time.Millisecond
		this.Server.ReadTimeout = time.Duration(this.opt.Option.TimeoutMilliseconds) * time.Millisecond
	}
	if proto == "https" {
		gls.Go(func() {
			this.Addr = U.FormatAddr(this.addr, proto)
			err := this.ListenAndServeTLS(gCertFile, gKeyFile)
			if err != nil {
				log.GetLog().Fatalf("Http.ServeTLS lis error:%s. addr:%s", err.Error(), this.addr)
			}
        })
    } else if proto == "http" {
		lis, err := net.Listen("tcp", U.FormatAddr(this.addr, proto))
		if err != nil {
			return err
		}

		gls.Go(func() {
			err := this.Serve(lis)
			if err != nil {
				log.GetLog().Fatalf("Http.Serve lis error:%s. addr:%s", err.Error(), this.addr)
			}
        })
    } else {
		return errors.Errorf("Error protocol:%s", proto)
    }

	this.served = true
	return nil
}
