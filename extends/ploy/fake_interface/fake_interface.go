package fake_interface

import (
	"fmt"
	"reflect"
	"golang.org/x/net/context"
	"github.com/yyzybb537/ketty"
)

type FakeInterface struct {
	m				map[string]*methodDesc
	realizedI		interface{}

}

type fakeFunc func(args ...interface{}) context.Context

type methodDesc struct {
	Name		string
	ArgsCount	int
	Handle		fakeFunc
}

func NewFakeInterface() *FakeInterface{
	return &FakeInterface{
		m:		make(map[string]*methodDesc),
    }
}

func (this *FakeInterface) Add(methodName string, argsCount int) error{
	if _, ok := this.m[methodName]; ok{
		return fmt.Errorf("conflict methodName: %s", methodName)
    }
	m := &methodDesc{
		Name:		methodName,
		ArgsCount:	argsCount,
	}
	this.m[methodName] = m
	return nil
}

func (this *FakeInterface) Realize(realizedPoint interface{}) error{
	if this.realizedI != nil {
		//已经实现就覆盖，如果覆盖出错也会清除已有的interface
		newFI := NewFakeInterface()
		this = newFI
	}
	for _, m := range this.m {
		handle, err := parse(realizedPoint, m.Name, m.ArgsCount)
		if err != nil {
			return err
        }
		m.Handle = handle
    }
	this.realizedI = realizedPoint
	return nil
}

func (this *FakeInterface) Interface() interface{} {
	return this.realizedI
}

func (this *FakeInterface) Do(methodName string, args ...interface{}) (ctx context.Context) {
	if this.realizedI == nil {
		return ketty.WithError(context.Background(), fmt.Errorf("invalid FakeInterface: not Realized"))
    }
	m, ok := this.m[methodName]
	if !ok {
		return ketty.WithError(ctx, fmt.Errorf("invalid MethodName: %s not Realized", methodName))
    }
	return m.Handle(args...)
}

func parse(realizedPoint interface{}, methodName string, argsCount int) (h fakeFunc, err error) {
	refType := reflect.TypeOf(realizedPoint)
	if refType.Kind() != reflect.Ptr {
		err = fmt.Errorf("invalid refType kind")
		return
	}

	refValue := reflect.ValueOf(realizedPoint);

	for i := 0; i < refType.NumMethod(); i++ {
		methodType := refType.Method(i)
		if methodType.Name != methodName {
			continue
		}

		if methodType.Type.NumIn() != argsCount + 1{
			err = fmt.Errorf("invalid argsCount: method %s, want: %d, real: %d", methodName, argsCount + 1, methodType.Type.NumIn())
			return
        }

		if methodType.Type.NumOut() != 1 {
			err = fmt.Errorf("invalid return argsCount: method %s, want: %d, real: %d", methodName, 1, methodType.Type.NumOut())
			return
        }

		if methodType.Type.Out(0).Name() != "Context" {
			err = fmt.Errorf("invalid return arg: method %s, want: Context, real: %s", methodName, 1, methodType.Type.Out(0).Name())
			return
        }

		methodValue := refValue.Method(i)
		h = func(args ...interface{}) (ctx context.Context){
			var in []reflect.Value
			for _, arg := range args {
				argValue := reflect.ValueOf(arg)
				in = append(in, argValue)
            }
			returnValues := methodValue.Call(in)
			var ok bool
			if ctx, ok = returnValues[0].Interface().(context.Context); !ok{
				return nil
			}
			return ctx
		}
		return
    }
	err = fmt.Errorf("no match method, want: %s", methodName)
	return
}
