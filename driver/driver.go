package driver

import (
	"fmt"
	"strings"
	U "github.com/yyzybb537/ketty/url"
)

type Driver interface {
	Watch(url U.Url) (up, down <-chan []U.Url, stop func(), err error)

	Register(url, value U.Url) (error)
}

var drivers = make(map[string]Driver)

func GetDriver(sDriver string) (Driver, error) {
	driver, exists := drivers[strings.ToLower(sDriver)]
	if !exists {
		return nil, fmt.Errorf("Unkown driver:%s", sDriver)
	}
	return driver, nil
}

func RegDriver(sDriver string, driver Driver) {
	drivers[strings.ToLower(sDriver)] = driver
}
