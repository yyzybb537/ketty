package ploy

import (
	//"github.com/yyzybb537/ketty/extends/log"
	//"github.com/yyzybb537/ketty/config"
	"github.com/yyzybb537/ketty/common"
	"github.com/yyzybb537/ketty"
)

type Server struct {
	sUrl		string
	sDriverUrl	string
	servers		[]ketty.Server
}

func NewServer(sUrl, sDriverUrl string) (server *Server, err error) {
	server = &Server{}
	server.sUrl = sUrl
	server.sDriverUrl = sDriverUrl
	return
}

func (this *Server) NewFlow(router string, handle common.ServiceHandle, implement interface{}) (flow FlowI, err error){
	server, err := ketty.Listen(this.sUrl + router, this.sDriverUrl)
	if err != nil {
		return
	}
	err = server.RegisterMethod(handle, implement)
	if err != nil {
		return
	}
	flow = NewBaseFlow()
	err = setInterface(implement, flow, "FlowI")
	if err != nil {
		return
	}

	this.servers = append(this.servers, server)
	return
}

func (this *Server) Serve() (err error){
	for _, s := range this.servers {
		err = s.Serve()
		if err != nil {
			return
		}
	}
	return
}
