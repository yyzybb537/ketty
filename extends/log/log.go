package log

import (
	"github.com/yyzybb537/ketty"
	"github.com/yyzybb537/ketty/log"
	"github.com/yyzybb537/ketty/extends/log/real_log"
)

func init() {
	el := &real_log.RealLog{}
	ketty.SetLog(el)
}

// ---------------------------------------------------
func Debugf(format string, args ... interface{}) {
	ketty.GetLog().Debugf(format, args ...)
}
func Infof(format string, args ... interface{}) {
	ketty.GetLog().Infof(format, args ...)
}
func Warningf(format string, args ... interface{}) {
	ketty.GetLog().Warningf(format, args ...)
}
func Errorf(format string, args ... interface{}) {
	ketty.GetLog().Errorf(format, args ...)
}
func Fatalf(format string, args ... interface{}) {
	ketty.GetLog().Fatalf(format, args ...)
}
func Recordf(format string, args ... interface{}) {
	ketty.GetLog().Recordf(format, args ...)
}
func Debugln(args ... interface{}) {
	ketty.GetLog().Debugln(args ...)
}
func Infoln(args ... interface{}) {
	ketty.GetLog().Infoln(args ...)
}
func Warningln(args ... interface{}) {
	ketty.GetLog().Warningln(args ...)
}
func Errorln(args ... interface{}) {
	ketty.GetLog().Errorln(args ...)
}
func Fatalln(args ... interface{}) {
	ketty.GetLog().Fatalln(args ...)
}
func Recordln(args ... interface{}) {
	ketty.GetLog().Recordln(args ...)
}
func S(key interface{}) log.LogI {
	return ketty.GetLog(key)
}
func EnableSection(key interface{}) {
	ketty.EnableLogSection(key)
}
func DisableSection(key interface{}) {
	ketty.DisableLogSection(key)
}
