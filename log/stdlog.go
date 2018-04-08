package log

import (
	"bytes"
	"fmt"
	"os"
)

type StdLog struct {
	opt *LogOption
}

func (this *StdLog) Clone(opt *LogOption) (LogI, error) {
	return &StdLog{opt : opt}, nil
}

func (this *StdLog) write(level Level, info string) {
	if this.opt == nil {
		this.opt = DefaultLogOption()
	}
	buf := bytes.NewBufferString("")
	this.opt.WriteHeader(level, 3, buf)
	buf.WriteString(info)
	fmt.Println(buf.String())
}

func (this *StdLog) logf(level Level, format string, args ... interface{}) {
	this.write(level, fmt.Sprintf(format, args...))
}
func (this *StdLog) logln(level Level, args ... interface{}) {
	this.write(level, fmt.Sprintln(args...))
}
func (this *StdLog) Debugf(format string, args ... interface{}) {
	this.logf(lv_debug, format, args...)
}
func (this *StdLog) Infof(format string, args ... interface{}) {
	this.logf(lv_info, format, args...)
}
func (this *StdLog) Warningf(format string, args ... interface{}) {
	this.logf(lv_warning, format, args...)
}
func (this *StdLog) Errorf(format string, args ... interface{}) {
	this.logf(lv_error, format, args...)
}
func (this *StdLog) Fatalf(format string, args ... interface{}) {
	this.logf(lv_fatal, format, args...)
}
func (this *StdLog) Recordf(format string, args ... interface{}) {
	this.logf(lv_record, format, args...)
}
func (this *StdLog) Debugln(args ... interface{}) {
	this.logln(lv_debug, args...)
}
func (this *StdLog) Infoln(args ... interface{}) {
	this.logln(lv_info, args...)
}
func (this *StdLog) Warningln(args ... interface{}) {
	this.logln(lv_warning, args...)
}
func (this *StdLog) Errorln(args ... interface{}) {
	this.logln(lv_error, args...)
}
func (this *StdLog) Fatalln(args ... interface{}) {
	this.logln(lv_fatal, args...)
}
func (this *StdLog) Recordln(args ... interface{}) {
	this.logln(lv_record, args...)
}
func (this *StdLog) Flush() error {
	return os.Stdout.Sync()
}
