package common

import (
	"fmt"
	"reflect"
)

type Manager struct {
	m map[string]interface{}
	handlerType interface{}
}

func NewManager(handlerType interface{}) *Manager {
	return &Manager {
		m : make(map[string]interface{}),
		handlerType : handlerType,
    }
}

func (this *Manager) Register(name string, obj interface{}) {
	ht := reflect.TypeOf(this.handlerType).Elem()
	st := reflect.TypeOf(obj)
	if !st.Implements(ht) {
		panic(fmt.Errorf("register %s. obj is not instance %s", name, ht.Name()))
	}
	this.register(name, obj)
}

func (this *Manager) register(name string, obj interface{}) {
	this.m[name] = obj
}

func (this *Manager) Get(name string) interface{} {
	obj, _ := this.m[name]
	return obj
}
