package oss2bdc

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config 配置
type Config struct {
	Oss     *ossConf
	SQL     *sqlConf
	ZipPath string
}

type ossConf struct {
	Endpoint   string
	Key        string
	Secret     string
	BucketName string
}

type sqlConf struct {
	Driver string
	Source string
}

// GetConfig 获取配置
func GetConfig() *Config {
	bytes, err := ioutil.ReadFile("./cfg/oss2bdc.json")

	if err != nil {
		log.Fatalln("ioutil.ReadFile: ", err)
	}

	config := new(Config)
	if err := json.Unmarshal(bytes, config); err != nil {
		log.Fatalln("json.Unmarshal: ", err)
	}

	return config
}
