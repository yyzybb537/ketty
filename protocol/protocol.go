package protocol

import (
	"fmt"
	"strings"
	U "github.com/yyzybb537/ketty/url"
)

type Protocol interface {
	CreateServer(url, driverUrl U.Url) (Server, error)

	Dial(url U.Url) (Client, error)
}

var protocols = make(map[string]Protocol)

func GetProtocol(sproto string) (Protocol, error) {
	proto, exists := protocols[strings.ToLower(sproto)]
	if !exists {
		return nil, fmt.Errorf("Unkown Protocol:%s", sproto)
	}
	return proto, nil
}

func RegProtocol(sproto string, proto Protocol) {
	protocols[strings.ToLower(sproto)] = proto
}
