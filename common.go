package ketty

import (
	"time"
)

func Assert(err error) {
	if err != nil {
		panic(err)
    }
}

func Hung() {
	time.Sleep(time.Second)
	GetLog().Infof("ketty service startup")
	for {
		time.Sleep(time.Second * 3600)
    }
}
