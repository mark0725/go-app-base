package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// 这里可以做一些设置和初始化操作
	initialDir, err := os.Getwd()
	if err != nil {
		// 处理错误
		panic(err)
	}

	newDir := filepath.Join(initialDir, "..")
	err = os.Chdir(newDir)
	if err != nil {
		// 处理错误
		panic(err)
	}
	defer func() {
		_ = os.Chdir(initialDir)
	}()

	// 运行所有测试，m.Run 返回一个状态码
	exitCode := m.Run()

	// 使用 m.Run 的返回值作为进程的退出码
	os.Exit(exitCode)
}

func Test_LoadConfig(t *testing.T) {
	config, err := LoadConfig(WithEnv("AI_CHAT"), WithConfigfile("./config.yaml"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(config)
}
