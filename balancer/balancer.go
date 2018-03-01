package ketty

import (
	"fmt"
	"strings"
	"golang.org/x/net/context"
	U "github.com/yyzybb537/ketty/url"
)

// @looks like grpc.Balancer
type Balancer interface {
	Filte(in []U.Url) (out []U.Url)

	Up(addr U.Url) (down func())

	Get(ctx context.Context) (addr U.Url, put func(), err error)

	Clone() Balancer
}

var balancers = make(map[string]Balancer)

func GetBalancer(sBalancer string) (Balancer, error) {
	balancer, exists := balancers[strings.ToLower(sBalancer)]
	if !exists {
		return nil, fmt.Errorf("Unkown balancer:%s", sBalancer)
	}
	return balancer.Clone(), nil
}

func RegBalancer(sBalancer string, balancer Balancer) {
	balancers[strings.ToLower(sBalancer)] = balancer
}
