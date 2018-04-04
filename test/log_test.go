package log_t

import (
	"testing"
	"github.com/yyzybb537/ketty"
	"github.com/yyzybb537/ketty/log"
	"time"
)

var _ = log.BindOption

func TestLog(t *testing.T) {
	ketty.GetLog().Debugf("stdout log")

	opt := log.DefaultLogOption()
	opt.HeaderFormat = "$L $datetime-ms $$ $gid $file:$line] "
	log.BindOption("key", opt)
	ketty.GetLog("key").Debugf("key log")

	opt.OutputFile = "logdir/ketty"
	flg, err := log.NewFileLog(opt)
	if err != nil {
		println(err.Error())
		return
	}
	log.SetLog(flg)
	ketty.GetLog().Debugf("file log")
	time.Sleep(time.Second)
}

func Benchmark_Log(b *testing.B) {
	b.StopTimer()
	opt := log.DefaultLogOption()
	opt.OutputFile = "logdir/ketty"
	flg, err := log.NewFileLog(opt)
	if err != nil {
		println(err.Error())
		return
	}
	log.SetLog(flg)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ketty.GetLog().Debugf("file log")
	}
//	log.FlushAll()
}

func Benchmark_LogNull(b *testing.B) {
	b.StopTimer()
	opt := log.DefaultLogOption()
	opt.OutputFile = "/dev/null"
	flg, err := log.NewFileLog(opt)
	if err != nil {
		println(err.Error())
		return
	}
	log.SetLog(flg)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ketty.GetLog().Debugf("file log")
	}
//	log.FlushAll()
}

type fakeWriter struct {}
func (this *fakeWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func Benchmark_Header(b *testing.B) {
	b.StopTimer()
	opt := log.DefaultLogOption()
	w := new(fakeWriter)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		opt.WriteHeader(log.Level(1), 0, w)
	}
}

