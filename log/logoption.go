package log

import (
	"io"
	"os"
	"fmt"
	"time"
	"strings"
	"bytes"
	"strconv"
	"runtime"
	"github.com/yyzybb537/gls"
	"unsafe"
	"context"
)

var _ unsafe.Pointer
var _ = strconv.IntSize

type LogOption struct {
	// toggle of logger
	Ignore bool

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
			OutputFile : "log/ketty",
			HeaderFormat : "$L $datetime-ms $gid $file:$line] ",
			RotateCategory : "size",
			RotateValue : 1912602624, // 1.8 GB
			RotateSuffixFormat : ".P$pid.$day-$hour-$second",
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
	var n int
	nBytes := 0
	ctx := context.Background()
	for _, cl := range this.callers {
		if str, ok := cl.(string); ok {
			_ = str
//			n, _ = buf.Write([]byte(str))
			n, _ = buf.Write(fastString2Bytes(str))
			nBytes += n
		} else if f, ok := cl.(func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context)); ok {
			_ = f
			n, ctx = f(ctx, level, depth, buf)
//			n, _ = buf.Write(fastString2Bytes(str))
			nBytes += n
		}
	}
	return nBytes
}

type lineKey struct{}
type fileKey struct{}
type timeKey struct{}

func headerLeader2func(leader string) func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context) {
	switch leader {
	case "level":
		return func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context) {
			n, _ := buf.Write(fastString2Bytes(level.ToString()))
			return n, ctx
		}

	case "L":
		return func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context) {
			n, _ := buf.Write([]byte(level.Header()))
			return n, ctx
		}

	case "datetime-ms":
		return func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context) {
			var tmp [27]byte
			t := time.Now()
			year, month, day := t.Date()
			hour, minute, second := t.Clock()
			usec := t.Nanosecond() / 1000
			tmp[0] = byte(year / 1000) + byte('0')
			tmp[1] = byte((year % 1000) / 100) + byte('0')
			tmp[2] = byte((year % 100) / 10) + byte('0')
			tmp[3] = byte(year % 10) + byte('0')
			tmp[4] = '-'
			tmp[5] = byte(month / 10) + byte('0')
			tmp[6] = byte(month % 10) + byte('0')
			tmp[7] = '-'
			tmp[8] = byte(day / 10) + byte('0')
			tmp[9] = byte(day % 10) + byte('0')
			tmp[10] = ' '
			tmp[11] = byte(hour / 10) + byte('0')
			tmp[12] = byte(hour % 10) + byte('0')
			tmp[13] = ':'
			tmp[14] = byte(minute / 10) + byte('0')
			tmp[15] = byte(minute % 10) + byte('0')
			tmp[16] = ':'
			tmp[17] = byte(second / 10) + byte('0')
			tmp[18] = byte(second % 10) + byte('0')
			tmp[19] = '.'
			tmp[20] = byte(usec / 100000) + byte('0')
			tmp[21] = byte((usec % 100000) / 10000) + byte('0')
			tmp[23] = byte((usec % 10000) / 1000) + byte('0')
			tmp[24] = byte((usec % 1000) / 100) + byte('0')
			tmp[25] = byte((usec % 100) / 10) + byte('0')
			tmp[26] = byte(usec % 10) + byte('0')
			n := 27
			buf.Write(tmp[:n])
			return n, ctx
		}

	case "datetime":
		return func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context) {
//			n, _ := buf.Write(fastString2Bytes(time.Now().Format("2006-01-02 15:04:05")))
			var tmp [19]byte
			t := time.Now()
			year, month, day := t.Date()
			hour, minute, second := t.Clock()
			tmp[0] = byte(year / 1000) + byte('0')
			tmp[1] = byte((year % 1000) / 100) + byte('0')
			tmp[2] = byte((year % 100) / 10) + byte('0')
			tmp[3] = byte(year % 10) + byte('0')
			tmp[4] = '-'
			tmp[5] = byte(month / 10) + byte('0')
			tmp[6] = byte(month % 10) + byte('0')
			tmp[7] = '-'
			tmp[8] = byte(day / 10) + byte('0')
			tmp[9] = byte(day % 10) + byte('0')
			tmp[10] = ' '
			tmp[11] = byte(hour / 10) + byte('0')
			tmp[12] = byte(hour % 10) + byte('0')
			tmp[13] = ':'
			tmp[14] = byte(minute / 10) + byte('0')
			tmp[15] = byte(minute % 10) + byte('0')
			tmp[16] = ':'
			tmp[17] = byte(second / 10) + byte('0')
			tmp[18] = byte(second % 10) + byte('0')
			n := 19
			buf.Write(tmp[:n])
			return n, ctx
		}

	case "date":
		return func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context) {
			var tmp [10]byte
			t := time.Now()
			year, month, day := t.Date()
			tmp[0] = byte(year / 1000) + byte('0')
			tmp[1] = byte((year % 1000) / 100) + byte('0')
			tmp[2] = byte((year % 100) / 10) + byte('0')
			tmp[3] = byte(year % 10) + byte('0')
			tmp[4] = '-'
			tmp[5] = byte(month / 10) + byte('0')
			tmp[6] = byte(month % 10) + byte('0')
			tmp[7] = '-'
			tmp[8] = byte(day / 10) + byte('0')
			tmp[9] = byte(day % 10) + byte('0')
			n := 10
			buf.Write(tmp[:n])
			return n, ctx
		}

	case "time":
		return func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context) {
			var tmp [8]byte
			t := time.Now()
			hour, minute, second := t.Clock()
			tmp[0] = byte(hour / 10) + byte('0')
			tmp[1] = byte(hour % 10) + byte('0')
			tmp[2] = ':'
			tmp[3] = byte(minute / 10) + byte('0')
			tmp[4] = byte(minute % 10) + byte('0')
			tmp[5] = ':'
			tmp[6] = byte(second / 10) + byte('0')
			tmp[7] = byte(second % 10) + byte('0')
			n := 8
			buf.Write(tmp[:n])
			return n, ctx
		}

	case "ms":
		return func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context) {
			var tmp [7]byte
			t := time.Now()
			usec := t.Nanosecond() / 1000
			tmp[0] = byte(usec / 100000) + byte('0')
			tmp[1] = byte((usec % 100000) / 10000) + byte('0')
			tmp[3] = byte((usec % 10000) / 1000) + byte('0')
			tmp[4] = byte((usec % 1000) / 100) + byte('0')
			tmp[5] = byte((usec % 100) / 10) + byte('0')
			tmp[6] = byte(usec % 10) + byte('0')
			n := 7
			buf.Write(tmp[:n])
			return n, ctx
		}

	case "gid":
		return func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context) {
			n, _ := buf.Write(fastString2Bytes(strconv.Itoa(int(gls.Goid()))))
			return n, ctx
		}

	case "file":
		return func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context) {
			file, ok := ctx.Value(fileKey{}).(string)
			if !ok {
				var line int
				file, line = getFileLine(3 + depth)
				ctx = context.WithValue(ctx, fileKey{}, file)
				ctx = context.WithValue(ctx, lineKey{}, line)
			}
			n, _ := buf.Write(fastString2Bytes(file))
			return n, ctx
		}

	case "line":
		return func(ctx context.Context, level Level, depth int, buf io.Writer) (int, context.Context) {
			line, ok := ctx.Value(lineKey{}).(int)
			if !ok {
				var file string
				file, line = getFileLine(3 + depth)
				ctx = context.WithValue(ctx, fileKey{}, file)
				ctx = context.WithValue(ctx, lineKey{}, line)
			}
			n, _ := buf.Write(fastString2Bytes(strconv.Itoa(line)))
			return n, ctx
		}

	default:
		return nil
	}
}

func getFileLine(depth int) (string, int) {
	_, file, line, ok := runtime.Caller(1 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return file, line
}

func fastString2Bytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
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

