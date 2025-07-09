package config

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

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

	rawconfig, err := func() (map[string]interface{}, error) {
		config := make(map[string]interface{})
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
	}()

	if err != nil {
		return nil, err
	}

	envVars := make(map[string]string)
	for _, envVar := range os.Environ() {
		// 分割"key=value"格式的字符串
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			//log.Trace("key:", key, "value:", value)
			envVars[key] = value
		}
	}

	data, _ := yaml.Marshal(rawconfig)
	values := map[string]interface{}{
		"info": getAppInfo(),
		"env":  envVars,
		"vars": rawconfig,
	}

	t := template.Must(template.New("template").Parse(string(data)))
	var buf bytes.Buffer
	if err := t.Execute(&buf, values); err != nil {
		log.Error("parse param error:", err)
		return nil, err
	}

	err = yaml.Unmarshal(buf.Bytes(), &config)
	if err != nil {
		log.Errorf("error: %v", err)
		return nil, err
	}

	return config, nil
}

type AppInfo struct {
	AppName      string
	AppRoot      string
	WorkDir      string
	AppVersion   string
	AppBuildTime string
	AppGitCommit string
}

type OSInfo struct {
	OS   string
	Arch string
}

func getAppInfo() map[string]interface{} {
	appinfo := AppInfo{
		AppRoot: getAppRoot(),
		WorkDir: getWorkDir(),
	}

	info := map[string]interface{}{
		"app": appinfo,
		"os": OSInfo{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
	}

	return info
}

func getAppRoot() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Trace("无法获取可执行文件路径: ", err)
		return ""
	}

	// 读取可执行文件的真实路径
	exeRealPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		log.Trace("无法解析符号链接: ", err)
		return ""
	}

	// 获取可执行文件所在的目录
	exeDir := filepath.Dir(exeRealPath)

	return exeDir

}

func getWorkDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Trace("Error getting current working directory: ", err)
		return ""
	}

	return wd
}
