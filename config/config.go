package config

import (
	"strings"
	"fmt"
)

func Read(cfg interface{}, file string) error {
	ext := getFileExtend(file)
	config, err := GetConfig(ext)
	if err != nil {
		return err
    }

	return config.Read(cfg, file)
}

func Write(cfg interface{}, file string) error {
	ext := getFileExtend(file)
	config, err := GetConfig(ext)
	if err != nil {
		return err
    }

	return config.Write(cfg, file)
}

// ------------------------- internal
type Config interface {
	Read(cfg interface{}, file string) error

	Write(cfg interface{}, file string) error
}

var configs = make(map[string]Config)

func GetConfig(sConfig string) (Config, error) {
	config, exists := configs[strings.ToLower(sConfig)]
	if !exists {
		return nil, fmt.Errorf("Unkown config file type:%s", sConfig)
	}
	return config, nil
}

func RegConfig(sConfig string, config Config) {
	configs[strings.ToLower(sConfig)] = config
}

func getFileExtend(file string) string {
	ss := strings.Split(file, ".")
	return ss[len(ss)-1]
}
