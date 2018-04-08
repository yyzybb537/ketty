package ploy

import (
	"github.com/yyzybb537/ketty"
	"github.com/yyzybb537/ketty/log"
	"fmt"
)

type Server struct {
	sUrl		string
	sDriverUrl	string
	servers		[]ketty.Server
	logKeys     []interface{}
}

func NewServer(sUrl, sDriverUrl string) (server *Server, err error) {
	server = &Server{}
	server.sUrl = sUrl
	server.sDriverUrl = sDriverUrl
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
	server, err := ketty.Listen(this.sUrl + router, this.sDriverUrl)
	if err != nil {
		return
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
	this.servers = append(this.servers, server)
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
			defer log.CleanupGlsDefaultKey()
		}
		err = s.Serve()
		if err != nil {
			return
		}
	}
	return
}
