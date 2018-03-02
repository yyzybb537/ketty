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
	for {
		time.Sleep(time.Second * 3600)
    }
}
