package log

import (
	"github.com/yyzybb537/ketty"
	"github.com/yyzybb537/ketty/log"
	"github.com/yyzybb537/glog"
)

type ExtendLog struct{}

func (this *ExtendLog) Debugf(format string, args ...interface{}) {
	//默认不打印
	glog.V(0).Infof(format, args)
}
func (this *ExtendLog) Infof(format string, args ...interface{}) {
	glog.Infof(format, args)
}
func (this *ExtendLog) Warningf(format string, args ...interface{}) {
	glog.Warningf(format, args)
}
func (this *ExtendLog) Errorf(format string, args ...interface{}) {
	glog.Errorf(format, args)
}
func (this *ExtendLog) Fatalf(format string, args ...interface{}) {
	glog.Fatalf(format, args)
}
func (this *ExtendLog) Recordf(format string, args ...interface{}) {
	//TODO 独立打印日志
	glog.Infof(format, args)
}
func (this *ExtendLog) Debugln(args ...interface{}) {
	glog.V(0).Infoln(args)
}
func (this *ExtendLog) Infoln(args ...interface{}) {
	glog.Infoln(args)
}
func (this *ExtendLog) Warningln(args ...interface{}) {
	glog.Warningln(args)
}
func (this *ExtendLog) Errorln(args ...interface{}) {
	glog.Errorln(args)
}
func (this *ExtendLog) Fatalln(args ...interface{}) {
	glog.Fatalln(args)
}
func (this *ExtendLog) Recordln(args ...interface{}) {
	//TODO 独立打印日志
	glog.Infoln(args)
}

// ---------------------------------------------------
func Debugf(format string, args ... interface{}) {
	ketty.GetLog().Debugf(format, args)
}
func Infof(format string, args ... interface{}) {
	ketty.GetLog().Infof(format, args)
}
func Warningf(format string, args ... interface{}) {
	ketty.GetLog().Warningf(format, args)
}
func Errorf(format string, args ... interface{}) {
	ketty.GetLog().Errorf(format, args)
}
func Fatalf(format string, args ... interface{}) {
	ketty.GetLog().Fatalf(format, args)
}
func Recordf(format string, args ... interface{}) {
	ketty.GetLog().Recordf(format, args)
}
func Debugln(args ... interface{}) {
	ketty.GetLog().Debugln(args)
}
func Infoln(args ... interface{}) {
	ketty.GetLog().Infoln(args)
}
func Warningln(args ... interface{}) {
	ketty.GetLog().Warningln(args)
}
func Errorln(args ... interface{}) {
	ketty.GetLog().Errorln(args)
}
func Fatalln(args ... interface{}) {
	ketty.GetLog().Fatalln(args)
}
func Recordln(args ... interface{}) {
	ketty.GetLog().Recordln(args)
}

func S(key interface{}) log.LogI {
	return ketty.GetLog(key)
}
