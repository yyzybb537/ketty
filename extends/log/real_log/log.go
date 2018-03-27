package real_log

import (
	"github.com/yyzybb537/glog"
)

type RealLog struct{}

func (this *RealLog) Debugf(format string, args ...interface{}) {
	//默认不打印
	glog.V(0).Infof(format, args ...)
}
func (this *RealLog) Infof(format string, args ...interface{}) {
	glog.Infof(format, args ...)
}
func (this *RealLog) Warningf(format string, args ...interface{}) {
	glog.Warningf(format, args ...)
}
func (this *RealLog) Errorf(format string, args ...interface{}) {
	glog.Errorf(format, args ...)
}
func (this *RealLog) Fatalf(format string, args ...interface{}) {
	glog.Fatalf(format, args ...)
}
func (this *RealLog) Recordf(format string, args ...interface{}) {
	//TODO 独立打印日志
	glog.Infof(format, args ...)
}
func (this *RealLog) Debugln(args ...interface{}) {
	glog.V(0).Infoln(args ...)
}
func (this *RealLog) Infoln(args ...interface{}) {
	glog.Infoln(args ...)
}
func (this *RealLog) Warningln(args ...interface{}) {
	glog.Warningln(args ...)
}
func (this *RealLog) Errorln(args ...interface{}) {
	glog.Errorln(args ...)
}
func (this *RealLog) Fatalln(args ...interface{}) {
	glog.Fatalln(args ...)
}
func (this *RealLog) Recordln(args ...interface{}) {
	//TODO 独立打印日志
	glog.Infoln(args ...)
}
