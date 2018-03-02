package auto

import (
	"flag"
)

var isAuto bool

func init() {
	flag.BoolVar(&isAuto, "auto", false, "Set ketty.auto.IsAuto flag")
	flag.Parse()
}

func IsAuto() bool {
	return isAuto
}
