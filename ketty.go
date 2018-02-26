package ketty

import (
)

// @主要组件
//type Protocol interface {}
//type Balancer interface {}
//type Client interface {}
//type Server interface {}

// @sUrl:   protocol://ip[:port][,ip[:port]]/path
// E.g:
//    http://127.0.0.1:8030/path
//    https://127.0.0.1:8030
//    tcp://127.0.0.1:8030
//    grpc://127.0.0.1:8030,127.0.0.1:8031
//
// @sBalanceUrl:  protocol://ip[:port][,ip[:port]]/path
//    etcd://127.0.0.1:2379/path
func Listen(sUrl, sDriverUrl string) (server Server, err error) {
	url, err := UrlFromString(sUrl)
	if err != nil {
		return
	}

	driverUrl, err := UrlFromString(sDriverUrl)
	if err != nil {
		return
	}

	proto, err := GetProtocol(url.Protocol)
	if err != nil {
		return
	}

	server, err = proto.CreateServer(url, driverUrl)
	server.AddAop(DefaultAop().GetAop()...)
	return
}

func Dial(sUrl, sBalancer string) (client Client, err error) {
	url, err := UrlFromString(sUrl)
	if err != nil {
		return
	}

	balancer, err := GetBalancer(sBalancer)
	if err != nil {
		return
	}

	clients := newClients(url, balancer, nil)
	err = clients.dial()
	if err != nil {
		return
    }

	client = clients
	client.AddAop(DefaultAop().GetAop()...)
	return
}

