package log

import (
	"fmt"
	"time"
	"runtime"
	"strings"
	"sync"
	"github.com/yyzybb537/gls"
)

// ---------------------------------------------------
type LogI interface {
	Debugf(format string, args ... interface{})
	Infof(format string, args ... interface{})
	Warningf(format string, args ... interface{})
	Errorf(format string, args ... interface{})
	Fatalf(format string, args ... interface{})
	Recordf(format string, args ... interface{})

	Debugln(args ... interface{})
	Infoln(args ... interface{})
	Warningln(args ... interface{})
	Errorln(args ... interface{})
	Fatalln(args ... interface{})
	Recordln(args ... interface{})
}

var logger LogI = &StdLog{}

func SetLog(l LogI) {
	logger = l
}

func GetLog(opt ... interface{}) LogI {
	if len(opt) == 0 {
		return logger
    }

	sectionMu.RLock()
	defer sectionMu.RUnlock()
	b, exists := sections[opt[0]]
	if !exists || b {
		return Verbose(true)
    }

	return Verbose(false)
}

// ---------------------------------------------------
type Level int
const (
	lv_debug Level = iota
	lv_info
	lv_warning
	lv_error
	lv_fatal
	lv_record
)

func (lv Level) ToString() string {
	switch lv {
	case lv_debug:
		return "Debug"
	case lv_info:
		return "Info"
	case lv_warning:
		return "Warning"
	case lv_error:
		return "Error"
	case lv_record:
		return "Record"
	case lv_fatal:
		return "Fatal"
	default:
		return "Unkown"
    }
}

func (lv Level) Header() string {
	switch lv {
	case lv_debug:
		return "D"
	case lv_info:
		return "I"
	case lv_warning:
		return "W"
	case lv_error:
		return "E"
	case lv_record:
		return "R"
	case lv_fatal:
		return "F"
	default:
		return "U"
    }
}

func now() string {
	return time.Now().Format("2006-01-02 15:04:05.999999")
}

func gid() int64 {
	return gls.Goid()
}

// ---------------------------------------------------
var sections = make(map[interface{}]bool)
var sectionMu sync.RWMutex

func EnableSection(opt interface{}) {
	sectionMu.Lock()
	defer sectionMu.Unlock()
	sections[opt] = true
}

func DisableSection(opt interface{}) {
	sectionMu.Lock()
	defer sectionMu.Unlock()
	sections[opt] = false
}

type Verbose bool
func (this Verbose) Debugf(format string, args ... interface{}) {
	if this {
		logger.Debugf(format, args ...)
    }
}
func (this Verbose) Infof(format string, args ... interface{}) {
	if this {
		logger.Infof(format, args ...)
    }
}
func (this Verbose) Warningf(format string, args ... interface{}) {
	if this {
		logger.Warningf(format, args ...)
    }
}
func (this Verbose) Errorf(format string, args ... interface{}) {
	if this {
		logger.Errorf(format, args ...)
    }
}
func (this Verbose) Fatalf(format string, args ... interface{}) {
	if this {
		logger.Fatalf(format, args ...)
    }
}
func (this Verbose) Recordf(format string, args ... interface{}) {
	if this {
		logger.Recordf(format, args ...)
    }
}
func (this Verbose) Debugln(args ... interface{}) {
	if this {
		logger.Debugln(args ...)
    }
}
func (this Verbose) Infoln(args ... interface{}) {
	if this {
		logger.Infoln(args ...)
    }
}
func (this Verbose) Warningln(args ... interface{}) {
	if this {
		logger.Warningln(args ...)
    }
}
func (this Verbose) Errorln(args ... interface{}) {
	if this {
		logger.Errorln(args ...)
    }
}
func (this Verbose) Fatalln(args ... interface{}) {
	if this {
		logger.Fatalln(args ...)
    }
}
func (this Verbose) Recordln(args ... interface{}) {
	if this {
		logger.Recordln(args ...)
    }
}
// ---------------------------------------------------


type FakeLog struct {}
func (this *FakeLog) Debugf(format string, args ... interface{}) {}
func (this *FakeLog) Infof(format string, args ... interface{}) {}
func (this *FakeLog) Warningf(format string, args ... interface{}) {}
func (this *FakeLog) Errorf(format string, args ... interface{}) {}
func (this *FakeLog) Fatalf(format string, args ... interface{}) {}
func (this *FakeLog) Recordf(format string, args ... interface{}) {}
func (this *FakeLog) Debugln(args ... interface{}) {}
func (this *FakeLog) Infoln(args ... interface{}) {}
func (this *FakeLog) Warningln(args ... interface{}) {}
func (this *FakeLog) Errorln(args ... interface{}) {}
func (this *FakeLog) Fatalln(args ... interface{}) {}
func (this *FakeLog) Recordln(args ... interface{}) {}

type StdLog struct {}
func (this *StdLog) header(level Level) string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	// header format: level time gid file:line]
	return fmt.Sprintf("%s %s %d %s:%d]", level.Header(), now(), gid(), file, line)
}
func (this *StdLog) logf(level Level, format string, args ... interface{}) {
	fmt.Printf("%s %s\n", this.header(level), fmt.Sprintf(format, args...))
}
func (this *StdLog) logln(level Level, args ... interface{}) {
	fmt.Println(this.header(level), fmt.Sprintln(args...))
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
