package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
	"fmt"
	"path/filepath"
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

	_, err = toml.DecodeFile(file, cfg)
	return
}

func (this *TomlConfig) Write(cfg interface{}, file string) (err error) {
	err = this.check(cfg)
	if err != nil {
		return
	}

	err = os.MkdirAll(filepath.Dir(file), 0775)
	if os.IsExist(err) {
		err = nil
	}
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
		err = fmt.Errorf("ConfigInit parameter must be a pointer. cfg.Name=%s cfg.Kind=%s", v.Type().Name(), v.Kind().String())
		return 
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		err = fmt.Errorf("ConfigInit parameter must be a pointer to struct. cfg.Name=%s cfg.Kind=%s", v.Type().Name(), v.Kind().String())
		return 
	}

	return 
}
