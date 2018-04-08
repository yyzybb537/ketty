package log_t

import (
	"testing"
	"github.com/yyzybb537/ketty"
	"github.com/yyzybb537/ketty/log"
	"time"
	"strings"
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

func Benchmark_FileLog(b *testing.B) {
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
//	opt.HeaderFormat = "$datetime"
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
}

func Benchmark_LogTime(b *testing.B) {
	b.StopTimer()
	w := new(fakeWriter)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}

func Benchmark_LogWithoutHeader(b *testing.B) {
	b.StopTimer()
	opt := log.DefaultLogOption()
	opt.HeaderFormat = ""
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
}

type fakeWriter struct {}
func (this *fakeWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

var _ = strings.Compare
func Benchmark_LogHeader(b *testing.B) {
	b.StopTimer()
	opt := log.DefaultLogOption()
//	opt.HeaderFormat = strings.Repeat("$datetime", 1)
//	opt.HeaderFormat = strings.Repeat("$gid", 1)
	w := new(fakeWriter)
	opt.WriteHeader(log.Level(1), 0, w)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		opt.WriteHeader(log.Level(1), 0, w)
	}
}

