package ploy

import (
	"github.com/yyzybb537/ketty"
	"github.com/yyzybb537/ketty/log"
	"github.com/yyzybb537/ketty/option"
	"fmt"
	"strings"
	"sync"
)

type Server struct {
	sUrl		string
	sDriverUrl	string
	servers		[]ketty.Server
	opt			option.OptionI
	logKeys     []interface{}
	ifGrpc		bool
	lock		sync.Mutex
}

func NewServer(sUrl, sDriverUrl string) (server *Server, err error) {
	server = &Server{}
	server.sUrl = sUrl
	server.sDriverUrl = sDriverUrl
	//不太优雅，但是ketty库目前没有实现返回server类型接口
	if strings.Contains(sUrl, "grpc") {
		server.ifGrpc = true
	}
	return
}

type getHandle interface{
	GetHandle() ketty.ServiceHandle
}

type getGlsLogKey interface{
	GlsLogKey() interface{}
}

func (this *Server) NewFlow(router string, implement interface{}) (flow FlowI, err error){
	h, ok := implement.(getHandle)
	if !ok {
		err = fmt.Errorf("not implement GetHandle")
		return
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	var server ketty.Server
	var needAppend bool
	if this.ifGrpc && len(this.servers) != 0 {
		server = this.servers[0]
	}else {
		server, err = ketty.Listen(this.sUrl + router, this.sDriverUrl)
		if err != nil {
			return
		}
		if this.opt != nil {
			err = server.SetOption(this.opt)
			if err != nil {
				return
			}
		}
		needAppend = true
	}
	
	err = server.RegisterMethod(h.GetHandle(), implement)
	if err != nil {
		return
	}
	flow = NewBaseFlow()
	err = setInterface(implement, flow, "FlowI")
	if err != nil {
		return
	}
	i, ok := implement.(_init)
	if ok {
		i.Init()
	}
	
	if needAppend {
		this.servers = append(this.servers, server)
	}
	this.appendLogKeys(implement)
	return
}

func (this *Server) NewOverLappedFlow(router string, implement interface{}, flow FlowI) (err error) {
	h, ok := implement.(getHandle)
	if !ok {
		err = fmt.Errorf("not implement getHandle")
		return
	}
	server, err := ketty.Listen(this.sUrl + router, this.sDriverUrl)
	if err != nil {
		return
	}
	err = server.RegisterMethod(h.GetHandle(), implement)
	if err != nil {
		return
	}
	err = setInterface(implement, flow, "FlowI")
	if err != nil {
		return
	}

	this.servers = append(this.servers, server)
	this.appendLogKeys(implement)
	return
}

func (this *Server) appendLogKeys(implement interface{}) {
	var logKey interface{}
	if logHandle, ok := implement.(getGlsLogKey); ok {
		logKey = logHandle.GlsLogKey()
	}
	this.logKeys = append(this.logKeys, logKey)
}

func (this *Server) Serve() (err error){
	for i, s := range this.servers {
		logKey := this.logKeys[i]
		if logKey != nil {
			log.SetGlsDefaultKey(logKey)
			defer log.CleanupGlsDefaultKey(logKey)
		}
		err = s.Serve()
		if err != nil {
			return
		}
	}
	return
}

func (this *Server) SetOption(opt option.OptionI) (err error){
	this.opt = opt
	for _, s := range this.servers {
		err = s.SetOption(this.opt)
		if err != nil {
			return
		}
	}
	return
}
