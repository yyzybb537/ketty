package log

import (
	"sync"
	"github.com/yyzybb537/gls"
	COM "github.com/yyzybb537/ketty/common"
	"github.com/pkg/errors"
)

// ---------------------------------------------------
type LogI interface {
	Clone(opt *LogOption) (LogI, error)

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

	Flush() error
}

var LogMgr *COM.Manager = COM.NewManager((*LogI)(nil))

var logger LogI
var logBindings = make(map[interface{}]LogI)
var logBindingsMtx sync.RWMutex

func init() {
	LogMgr.Register("std", new(StdLog))
	LogMgr.Register("file", new(FileLog))
	LogMgr.Register("fake", new(FakeLog))

	if logger == nil {
		var err error
		opt := DefaultLogOption()
		logger, err = MakeLogger(opt)
		if err != nil {
			opt.LogCategory = "std"
			logger, err = MakeLogger(opt)
			if err != nil {
				panic(err)
			}
		}
	}
}

func MakeLogger(opt *LogOption) (LogI, error) {
	lg, ok := LogMgr.Get(opt.LogCategory).(LogI)
	if !ok {
		return nil, errors.Errorf("Unkown log category in option. LogCategory=%s", opt.LogCategory)
	}

	newLg, err := lg.Clone(opt)
	if err != nil {
		return nil, err
	}

	return newLg, nil
}

func SetLog(l LogI) {
	if logger != l {
		FlushAll()
		logBindingsMtx.Lock()
		defer logBindingsMtx.Unlock()
		logBindings = make(map[interface{}]LogI)
		logger = l
	}
}

func BindOption(key interface{}, opt *LogOption) (LogI, error) {
	logBindingsMtx.Lock()
	defer logBindingsMtx.Unlock()
	lg, err := MakeLogger(opt)
	if err != nil {
		return nil, err
	}
	logBindings[key] = lg
	return lg, nil
}

func GetLog(keys ... interface{}) LogI {
	var key interface{}
	if len(keys) > 0 {
		key = keys[0]
    } else {
		key = GetGlsDefaultKey()
		if key == nil {
			return logger
		}
	}

	if !CheckSection(key) {
		return gFakeLog
	}

	logBindingsMtx.RLock()
	defer logBindingsMtx.RUnlock()
	if lg, exists := logBindings[key]; exists {
		return lg
	}
	return logger
}

func FlushAll() {
	logBindingsMtx.RLock()
	defer logBindingsMtx.RUnlock()
	for _, lg := range logBindings {
		lg.Flush()
	}
	logger.Flush()
}

type priGlsKey struct{}

// 警告: 调用这个接口的地方, 一定要记得使用结束后调用Cleanup, 否则会内存泄漏
func SetGlsDefaultKey(key interface{}) {
	gls.Set(priGlsKey{}, key)
}

func CleanupGlsDefaultKey(key interface{}) {
	glsKey := GetGlsDefaultKey()
	if glsKey == key {
		gls.Del(priGlsKey{})
	}
}

func GetGlsDefaultKey() interface{} {
	return gls.Get(priGlsKey{})
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
