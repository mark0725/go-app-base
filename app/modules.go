package app

import (
	"container/list"
	"errors"
	"fmt"
)

type IAppModule interface {
	Init(appConfig interface{}, depends []string) error
}

type AppModule struct {
	Name    string
	Module  IAppModule
	Depends []string
	//Done    chan Signal
}

var g_appModules = make(map[string]AppModule, 10)

func AppModuleRegister(name string, module IAppModule, depends []string) {
	g_appModules[name] = AppModule{
		Name:    name,
		Module:  module,
		Depends: depends,
		//Done:    make(chan Signal),
	}
}

func InitializeModules(appConfig interface{}) error {
	// 保存每个模块被依赖的次数
	inDegree := make(map[string]int)

	// 构建依赖图
	for name, module := range g_appModules {
		if _, ok := inDegree[name]; !ok {
			inDegree[name] = 0
		}

		for _, dep := range module.Depends {
			if _, ok := inDegree[dep]; !ok {
				inDegree[dep] = 1
			} else {
				inDegree[dep]++
			}
		}
	}
	//fmt.Println("inDegree: ", inDegree)
	// 使用队列进行拓扑排序
	queue := list.New()
	for name, degree := range inDegree {
		if degree == 0 {
			queue.PushBack(name)
		}
	}

	var sortedModules []string
	for queue.Len() > 0 {

		element := queue.Front()
		queue.Remove(element)
		name := element.Value.(string)
		sortedModules = append([]string{name}, sortedModules...)

		module, ok := g_appModules[name]
		if !ok {
			return fmt.Errorf("module %s not found", name)
		}

		for _, dep := range module.Depends {
			if inDegree[dep] > 0 {
				inDegree[dep]--
				if inDegree[dep] == 0 {
					queue.PushBack(dep)
				}
			}
		}

		//fmt.Println("sortedModules: ", sortedModules)
	}

	// 检查是否有环
	if len(sortedModules) != len(g_appModules) {
		fmt.Println("sortedModules: ", sortedModules)
		fmt.Println("appModules: ", g_appModules)
		return errors.New("circular dependency detected")
	}

	// 按顺序初始化模块
	for _, name := range sortedModules {
		module := g_appModules[name]
		err := module.Module.Init(appConfig, module.Depends)
		if err != nil {
			return fmt.Errorf("failed to initialize module %s: %w", name, err)
		}
	}

	return nil
}
