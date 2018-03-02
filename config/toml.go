package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
	"fmt"
)

type TomlConfig struct {}

func init() {
	RegConfig("toml", new(TomlConfig))
}

func (this *TomlConfig) Read(cfg interface{}, file string) (err error) {
	err = this.check(cfg)
	if err != nil {
		return
	}

	_, err = toml.DecodeFile(os.Args[1], cfg)
	return
}

func (this *TomlConfig) Write(cfg interface{}, file string) (err error) {
	err = this.check(cfg)
	if err != nil {
		return
	}

	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()
	err = toml.NewEncoder(f).Encode(cfg)
	return
}

func (this *TomlConfig) check(cfg interface{}) (err error) {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr {
		err = fmt.Errorf("ConfigInit parameter must be a pointer")
		return 
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		err = fmt.Errorf("ConfigInit parameter must be a pointer to struct")
		return 
	}

	return 
}
