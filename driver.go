package ketty

import (
	"fmt"
	"strings"
)

type Driver interface {
	DefaultPort() int

	Watch(url Url) (up, down <-chan []Url, stop func(), err error)

	Register(url, value Url) (error)
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
