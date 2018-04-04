package log

import (
	"io"
	"os"
	"fmt"
	"time"
	"strings"
	"bytes"
//	"strconv"
	"runtime"
	"github.com/yyzybb537/gls"
	"unsafe"
)

type LogOption struct {
	// toggle of logger
	Active bool

	// output file name
	OutputFile string

	// log info format: "$Level $datetime-ms $gid $file:$line] "
	HeaderFormat string

	// rotate category: size / time / ""
	RotateCategory string

	// rotate parameter: bytes of one file / seconds
	RotateValue int64

	// rotate suffix format of filename: .$pid.$day-$hour-$second.log
	RotateSuffixFormat string

	// extends fields
	Extends map[string]string

	hf *formater
	sf *formater
}

func (this *LogOption) Clone() *LogOption {
	opt := new(LogOption)
	*opt = *this
	opt.hf = nil
	opt.sf = nil
	opt.Extends = make(map[string]string)
	for k, v := range this.Extends {
		opt.Extends[k] = v
	}
	return opt
}

var defaultOpt *LogOption
func DefaultLogOption() *LogOption {
	if defaultOpt == nil {
		defaultOpt = &LogOption{
			Active : true,
			OutputFile : "ketty",
			HeaderFormat : "$L $date $time$ms $gid $file:$line] ",
			RotateCategory : "size",
			RotateValue : 1912602624, // 1.8 GB
			RotateSuffixFormat : ".P$pid.$day-$hour-$second.log",
		}
	}

	return defaultOpt.Clone()
}

type formater struct {
	callers []interface{}
}

func splitBefore(s string, sep string) []string {
	ss := strings.Split(s, sep)
	for i := 1; i < len(ss); i++ {
		ss[i] = sep + ss[i]
	}
	return ss
}

func createFormater(format string, conv func(string)(interface{}, error)) *formater {
	hf := &formater{}
	sss := strings.Split(format, "$$")
	for idx, str := range sss {
		ss := splitBefore(str, "$")
		for _, s := range ss {
			if s == "" {
				continue
			}

			if s[0] != '$' {
				hf.addCaller(s)
				continue
			}

			parsed := false
			for i := 0; i < len(s) - 1; i++ {
				leader := s[1:len(s) - i]
				f, err := conv(leader)
				if err == nil {
					hf.addCaller(f)
					if i > 0 {
						hf.addCaller(s[len(s) - i:])
					}
					parsed = true
					break
				}
			}
			if parsed {
				continue
			}

			hf.addCaller(s)
		}

		if idx + 1 < len(sss) {
			hf.addCaller("$")
		}
	}
	hf.merge()
	return hf
}

func (this *formater) addCaller(arg interface{}) {
//	println("addCaller", arg)
	this.callers = append(this.callers, arg)
}

func (this *formater) merge() {
	newCallers := []interface{}{}
	for _, caller := range this.callers {
		if len(newCallers) == 0 {
			newCallers = append(newCallers, caller)
			continue
		}

		if c, ok := caller.(string); ok {
			if last, ok := newCallers[len(newCallers)-1].(string); ok {
				newCallers[len(newCallers)-1] = c + last
				continue
			}
		}

		newCallers = append(newCallers, caller)
	}
	this.callers = newCallers
}

func (this *formater) WriteHeader(level Level, depth int, buf io.Writer) int {
	nBytes := 0
	for _, cl := range this.callers {
		if str, ok := cl.(string); ok {
//			n, _ := fmt.Fprint(buf, str)
//			n, _ := buf.Write([]byte(str))
			n, _ := buf.Write(*(*[]byte)(unsafe.Pointer(&str)))
			nBytes += n
		} else if f, ok := cl.(func(level Level, depth int, buf io.Writer) int); ok {
//			n, _ := fmt.Fprint(buf, f(level, depth))
//			n, _ := buf.Write([]byte(f(level, depth)))
			n := f(level, depth, buf)
//			n, _ := buf.Write(*(*[]byte)(unsafe.Pointer(&str)))
			nBytes += n
		}
	}
	return nBytes
}

func headerLeader2func(leader string) func(level Level, depth int, buf io.Writer) int {
	switch leader {
	case "level":
		return func(level Level, depth int, buf io.Writer) int {
			n, _ := fmt.Fprint(buf, level.ToString())
			return n
		}

	case "L":
		return func(level Level, depth int, buf io.Writer) int {
			n, _ := fmt.Fprint(buf, level.Header())
			return n
		}

	case "datetime-ms":
		return func(level Level, depth int, buf io.Writer) int {
			n, _ := fmt.Fprint(buf, time.Now().Format("2006-01-02 15:04:05.999999"))
			return n
		}

	case "datetime":
		return func(level Level, depth int, buf io.Writer) int {
			n, _ := fmt.Fprint(buf, time.Now().Format("2006-01-02 15:04:05"))
			return n
		}

	case "date":
		return func(level Level, depth int, buf io.Writer) int {
			n, _ := fmt.Fprint(buf, time.Now().Format("2006-01-02"))
			return n
		}

	case "time":
		return func(level Level, depth int, buf io.Writer) int {
			n, _ := fmt.Fprint(buf, time.Now().Format("15:04:05"))
			return n
		}

	case "ms":
		return func(level Level, depth int, buf io.Writer) int {
			n, _ := fmt.Fprint(buf, time.Now().Format(".999999"))
			return n
		}

	case "gid":
		return func(level Level, depth int, buf io.Writer) int {
			n, _ := fmt.Fprintf(buf, "%d", gls.Goid())
			return n
		}

	case "file":
		return func(level Level, depth int, buf io.Writer) int {
			_, file, _, ok := runtime.Caller(3 + depth)
			if !ok {
				file = "???"
			} else {
				slash := strings.LastIndex(file, "/")
				if slash >= 0 {
					file = file[slash+1:]
				}
			}
			n, _ := fmt.Fprint(buf, file)
			return n
		}

	case "line":
		return func(level Level, depth int, buf io.Writer) int {
			_, _, line, ok := runtime.Caller(3 + depth)
			if !ok {
				line = 1
			}
			n, _ := fmt.Fprintf(buf, "%d", line)
			return n
		}

	default:
		return nil
	}
}

func suffixLeader2func(leader string) func() string {
	switch leader {
	case "datetime":
		return func()string {
			return time.Now().Format("2006-01-02T15:04:05")
		}

	case "date":
		return func()string {
			return time.Now().Format("2006-01-02")
		}

	case "time":
		return func()string {
			return time.Now().Format("15-04-05")
		}

	case "day":
		return func()string {
			return fmt.Sprintf("%02d", time.Now().Day())
		}

	case "hour":
		return func()string {
			return fmt.Sprintf("%02d", time.Now().Hour())
		}

	case "minute":
		return func()string {
			return fmt.Sprintf("%02d", time.Now().Minute())
		}

	case "second":
		return func()string {
			return fmt.Sprintf("%02d", time.Now().Second())
		}

	case "pid":
		return func()string {
			return fmt.Sprintf("%02d", os.Getpid())
		}

	default:
		return nil
	}
}

func (this *formater) GetSuffix() string {
	buf := bytes.NewBufferString("")
	for _, cl := range this.callers {
		if str, ok := cl.(string); ok {
			buf.Write([]byte(str))
		} else if f, ok := cl.(func() string); ok {
			buf.Write([]byte(f()))
		}
	}
	return buf.String()
}

func (this *LogOption) WriteHeader(level Level, depth int, buf io.Writer) int {
	if this.hf == nil {
		this.hf = createFormater(this.HeaderFormat, func(leader string)(interface{}, error){
			f := headerLeader2func(leader)
			if f != nil {
				return f, nil
			}
			return nil, fmt.Errorf("error")
		})
	}

	return this.hf.WriteHeader(level, depth, buf)
}

func (this *LogOption) GetSuffix() string {
	if this.sf == nil {
		this.sf = createFormater(this.RotateSuffixFormat, func(leader string)(interface{}, error){
			f := suffixLeader2func(leader)
			if f != nil {
				return f, nil
			}
			return nil, fmt.Errorf("error")
		})
	}

	return this.sf.GetSuffix()
}

