package ketty

import (
	"fmt"
	"strings"
	"golang.org/x/net/context"
	
)

// @looks like grpc.Balancer
type Balancer interface {
	Filte(in []Url) (out []Url)

	Up(addr Url) (down func())

	Get(ctx context.Context) (addr Url, put func(), err error)

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
