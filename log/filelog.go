package log

import (
	"fmt"
	"os"
	"io"
	"github.com/pkg/errors"
	"sync"
	"bufio"
	"time"
	"path/filepath"
)

const cBufSize = 256 * 1024
const cFlushCycle time.Duration = time.Millisecond * 100

type FileLog struct {
	opt *LogOption
	f syncWriter
	w *bufio.Writer
	mu sync.Mutex
	nWriteBytes int64
}

type syncWriter interface {
	io.Writer
	Sync() error
}

func NewFileLog(opt *LogOption) (*FileLog, error) {
	filelg := &FileLog{
		opt : opt,
	}
	if err := filelg.reopen(); err != nil {
		return nil, err
	}
	filelg.goHouseKeeper()
	return filelg, nil
}

func (this *FileLog) Clone(opt *LogOption) (LogI, error) {
	return NewFileLog(opt)
}

//func (this *FileLog) write(level Level, format string, args ... interface{}) {
func (this *FileLog) write(level Level, info string) {
	if this.opt.Ignore {
		return
	}

	this.mu.Lock()
	defer this.mu.Unlock()
	n := this.opt.WriteHeader(level, 3, this.w)
	this.nWriteBytes += int64(n)
	/*if format != "" {
		n, _ = fmt.Fprintf(this.w, format, args...)
	}else {
		n, _ = fmt.Fprintln(this.w, args...)
	}*/
	n, _ = this.w.WriteString(info)
	this.nWriteBytes += int64(n)
	this.w.WriteByte('\n')
	this.nWriteBytes += 1

	if this.opt.RotateCategory == "size" {
		if this.nWriteBytes >= this.opt.RotateValue {
			this.nWriteBytes = 0
			this.rotateWithoutLock()
		}
	}
}

func (this *FileLog) reopen() error {
	var err error
	_ = err
	if this.w != nil {
		this.w.Flush()
	}
	if this.f != nil {
		this.f.Sync()
	}
	this.nWriteBytes = 0

	if this.opt.OutputFile == "stdout" {
		this.f = os.Stdout
	} else if this.opt.OutputFile == "stderr" {
		this.f = os.Stderr
	} else if this.opt.OutputFile == os.DevNull {
		this.f, err = os.OpenFile(this.opt.OutputFile, os.O_WRONLY, 0666)
		if err != nil {
			return errors.Errorf("open log file(%s) error:%s\n", this.opt.OutputFile, err.Error())
		}
	} else {
		filename := this.opt.OutputFile + this.opt.GetSuffix() + ".log"
		f, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(filename), 0775)
		}
		if err == nil {
			f, err = os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		}
		if err != nil {
			return errors.Errorf("open log file(%s) error:%s\n", this.opt.OutputFile, err.Error())
		}

		fileInfo, err := f.Stat()
		if err != nil {
			return errors.Errorf("open log file.Stat(%s) error:%s\n", this.opt.OutputFile, err.Error())
		}
		this.nWriteBytes = fileInfo.Size()
		this.f = f

		symlink := this.opt.OutputFile + ".log"
		if symlink != filename {
			os.Remove(symlink)
			os.Symlink(filepath.Base(filename), symlink)
		}
	}

	this.w = bufio.NewWriterSize(this.f, cBufSize)

	return nil
}

func (this *FileLog) rotate() {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.rotateWithoutLock()
}

func (this *FileLog) rotateWithoutLock() {
	this.reopen()
}

func (this *FileLog) goHouseKeeper() {
	go func() {
		for {
			time.Sleep(cFlushCycle)

			this.mu.Lock()
			this.w.Flush()
			this.mu.Unlock()

			this.f.Sync()
		}
	}()

	if this.opt.RotateCategory == "time" {
		go func() {
			for {
				time.Sleep(time.Duration(this.opt.RotateValue) * time.Second)
				this.rotate()
			}
		}()
	}
}

func (this *FileLog) logf(level Level, format string, args ... interface{}) {
	//this.write(level, format, args...)
	this.write(level, fmt.Sprintf(format, args...))
}
func (this *FileLog) logln(level Level, args ... interface{}) {
	//this.write(level, "", args...)
	this.write(level, fmt.Sprintln(args...))
}
func (this *FileLog) Debugf(format string, args ... interface{}) {
	this.logf(lv_debug, format, args...)
}
func (this *FileLog) Infof(format string, args ... interface{}) {
	this.logf(lv_info, format, args...)
}
func (this *FileLog) Warningf(format string, args ... interface{}) {
	this.logf(lv_warning, format, args...)
}
func (this *FileLog) Errorf(format string, args ... interface{}) {
	this.logf(lv_error, format, args...)
}
func (this *FileLog) Fatalf(format string, args ... interface{}) {
	this.logf(lv_fatal, format, args...)
}
func (this *FileLog) Recordf(format string, args ... interface{}) {
	this.logf(lv_record, format, args...)
}
func (this *FileLog) Debugln(args ... interface{}) {
	this.logln(lv_debug, args...)
}
func (this *FileLog) Infoln(args ... interface{}) {
	this.logln(lv_info, args...)
}
func (this *FileLog) Warningln(args ... interface{}) {
	this.logln(lv_warning, args...)
}
func (this *FileLog) Errorln(args ... interface{}) {
	this.logln(lv_error, args...)
}
func (this *FileLog) Fatalln(args ... interface{}) {
	this.logln(lv_fatal, args...)
}
func (this *FileLog) Recordln(args ... interface{}) {
	this.logln(lv_record, args...)
}
func (this *FileLog) Flush() error {
	this.mu.Lock()
	if err := this.w.Flush(); err != nil {
		this.mu.Unlock()
		return err
	}
	this.mu.Unlock()
	return this.f.Sync()
}
