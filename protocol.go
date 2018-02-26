package ketty

import (
	"fmt"
	"strings"
)

type Protocol interface {
	DefaultPort() int

	CreateServer(url, driverUrl Url) (Server, error)

	Dial(url Url) (Client, error)
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
