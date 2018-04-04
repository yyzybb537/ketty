package log

import (
	"sync"
)

var sections = make(map[interface{}]bool)
var sectionMu sync.RWMutex

func CheckSection(opt interface{}) bool {
	if len(sections) == 0 {
		return true
	}

	sectionMu.RLock()
	defer sectionMu.RUnlock()
	if ok, exists := sections[opt]; exists {
		return ok
	}
	return true
}

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

type FakeLog struct {}
func (this *FakeLog) Clone(opt *LogOption) (LogI, error) { return &FakeLog{}, nil }
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
func (this *FakeLog) Flush() error { return nil }

var gFakeLog = &FakeLog{}
