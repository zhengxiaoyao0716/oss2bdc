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
	RawPath string
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

var config *Config

func init() {
	bytes, err := ioutil.ReadFile("./cfg/oss2bdc.json")

	if err != nil {
		log.Fatalln("ioutil.ReadFile: ", err)
	}

	config = new(Config)
	if err := json.Unmarshal(bytes, config); err != nil {
		log.Fatalln("json.Unmarshal: ", err)
	}
}

// GetConfig 获取配置
func GetConfig() *Config {
	return config
}
