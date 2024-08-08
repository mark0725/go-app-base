package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	env2toml "github.com/mark0725/env2toml-go"
	"github.com/mark0725/go-app-base/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func loadDotEnv(prefix string) error {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	envname := os.Getenv(prefix + "_ENV")
	if envname == "" {
		return nil
	}

	envFilename := ".env." + envname
	err = godotenv.Load(envFilename)
	if err != nil {
		log.Error("Error loading file", envFilename)
		return err
	}

	return nil

}

func loadConfigFile(path string) (map[string]interface{}, error) {
	// 读取 YAML 文件
	data, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("error: %v", err)
		return nil, err
	}

	// 创建一个map
	var config map[string]interface{}

	// 解析 YAML 数据到map
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Errorf("error: %v", err)
		return nil, err
	}

	// 打印解析的结果
	//log.Tracef("yaml config: %#v\n", config)
	return config, nil
}

func LoadConfig(opts ...ConfigOption) (map[string]interface{}, error) {
	options := buildOptions(opts...)
	config := make(map[string]interface{})

	if options == nil {
		return config, nil
	}

	if options.configFile != "" {
		yamlConfg, err := loadConfigFile(options.configFile)
		if err != nil {
			log.Error("Error:", err)
			return nil, err
		}

		log.Tracef("config file data: %#v\n", yamlConfg)
		config = utils.MergeMaps(config, yamlConfg)

	}

	if options.prefix != "" {
		var envConfigs map[string]interface{} = nil
		_ = loadDotEnv(options.prefix)

		tomlData, err := env2toml.Parse(options.prefix + "_")
		if err != nil {
			log.Error("Error:", err)
			return nil, err
		}
		log.Trace(tomlData)

		_, err = toml.Decode(tomlData, &envConfigs)
		if err != nil {
			log.Error("parse toml error", err)
			return nil, err
		}

		log.Tracef("env config: %#v\n", envConfigs)

		config = utils.MergeMaps(config, envConfigs)
	}

	return config, nil
}
