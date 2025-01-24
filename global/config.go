package global

import (
	"os"
	"scibe/utils/config"

	"gopkg.in/yaml.v3"
)

var cfg config.Config

func InitConfig(fp string) {
	file, err := os.Open(fp)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		panic(err)
	}
}

func Config() config.Config {
	return cfg
}
