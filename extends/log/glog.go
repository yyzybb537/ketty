package log

import (
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
