package ketty

import (
	U "github.com/yyzybb537/ketty/url"
	P "github.com/yyzybb537/ketty/protocol"
	B "github.com/yyzybb537/ketty/balancer"
	A "github.com/yyzybb537/ketty/aop"
	_ "github.com/yyzybb537/ketty/protocol/grpc"
	_ "github.com/yyzybb537/ketty/protocol/http"
)

type Dummy interface{}
type Client P.Client
type Server P.Server

// @目前支持的组件
//1.Protocol
//   grpc, http, https
//2.Balancer
//   robin
//3.Driver for find service
//   etcd
//4.AOP
//   cost, exception, logger, trace

// @后续计划
//1.Protocol
//   none
//2.Balancer
//   conhash, random
//3.Driver for find service
//   zookeeper
//* Driver mode:
//   only-master, master-groups
//4.AOP
//   timeout, trace-timeout, statistics
//5.Log
//   glog or seelog

// @sUrl:   protocol://ip[:port][,ip[:port]]/path
//  E.g:
//    http://127.0.0.1:8030/path
//    https://127.0.0.1:8030
//    grpc://127.0.0.1:8030,127.0.0.1:8031
//    
//  Support marshal: pb(default), json.
//  Support http method: GET, POST(default)
//  Support http data transport: body(default), query, multipart
//
//  User can make extend marshal.
//  If use `query`, the marshal will be ignored.
//
//  Others E.g:
//    http.json://127.0.0.1:8030/path
//    http.json.query://127.0.0.1:8030/path
//    http.pb.query://127.0.0.1:8030/path
//    http.json.get.query://127.0.0.1:8030/path
//    http.json.post.multipart://127.0.0.1:8030/path
//    http.get.query://127.0.0.1:8030/path
//    http.post.body://127.0.0.1:8030/path
//
// @sBalanceUrl:  driver://ip[:port][,ip[:port]]/path
//    etcd://127.0.0.1:2379/path
func Listen(sUrl, sDriverUrl string) (server Server, err error) {
	url, err := U.UrlFromString(sUrl)
	if err != nil {
		return
	}

	driverUrl, err := U.UrlFromString(sDriverUrl)
	if err != nil {
		return
	}

	proto, err := P.GetProtocol(url.GetMainProtocol())
	if err != nil {
		return
	}

	server, err = proto.CreateServer(url, driverUrl)
	server.AddAop(A.DefaultAop().GetAop()...)
	return
}

func Dial(sUrl, sBalancer string) (client Client, err error) {
	url, err := U.UrlFromString(sUrl)
	if err != nil {
		return
	}

	balancer, err := B.GetBalancer(sBalancer)
	if err != nil {
		return
	}

	clients := newClients(url, balancer, nil)
	err = clients.dial()
	if err != nil {
		return
    }

	client = clients
	client.AddAop(A.DefaultAop().GetAop()...)
	return
}

