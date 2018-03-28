package ploy

import (
	"reflect"
	"fmt"
)

func setInterface(owner interface{}, impl interface{}, interfaceName string) (err error){
	refV := reflect.ValueOf(owner)
	if refV.Kind() != reflect.Ptr {
		err = fmt.Errorf("invalid refType kind")
		return
	}
	refV = refV.Elem()
	for i := 0;i < refV.NumField(); i++ {
		field := refV.Type().Field(i)
		if field.Anonymous && field.Name == interfaceName {
			refV.Field(i).Set(reflect.ValueOf(impl))
			return
		}
	}
	err = fmt.Errorf("not such interfaceName: %s", interfaceName)
	return
}
