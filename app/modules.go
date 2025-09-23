package app

import (
	"container/list"
	"errors"
	"fmt"

	base_utils "github.com/mark0725/go-app-base/utils"
)

type IAppModule interface {
	Init(appConfig interface{}, depends []string) error
}

type AppModule struct {
	Name    string
	Module  IAppModule
	Depends []string
	Options AppModuleRegisterOptions
	Ready   bool
	//Done    chan Signal
}

type AppModuleOptions struct {
	modules []string
}

var g_appModules = make(map[string]*AppModule, 0)

type AppModuleOption func(*AppModuleOptions)

func WithModules(modules []string) AppModuleOption {
	return func(m *AppModuleOptions) {
		m.modules = modules
	}
}

func GetModules() map[string]*AppModule {
	return g_appModules
}

func GetReadyModules() []string {
	modules := make([]string, 0)
	for _, module := range g_appModules {
		if module.Ready {
			modules = append(modules, module.Name)
		}
	}
	return modules
}

func GetModule(name string) (AppModule, bool) {
	module, ok := g_appModules[name]
	return *module, ok
}

type AppModuleRegisterOptions struct {
	config any
}
type AppModuleRegisterOption func(*AppModuleRegisterOptions)

func AppModuleRegisterOptionWithConfigType(config any) AppModuleRegisterOption {
	return func(o *AppModuleRegisterOptions) {
		o.config = config
	}
}

func AppModuleRegister(name string, module IAppModule, depends []string, opts ...AppModuleRegisterOption) {
	options := AppModuleRegisterOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	g_appModules[name] = &AppModule{
		Name:    name,
		Module:  module,
		Depends: depends,
		Options: options,
		//Done:    make(chan Signal),
	}
}

func loadModuleConfig(config any) error {
	return base_utils.MapToStruct(g_AppConfigOrig, config)
}

func InitializeModules(appConfig interface{}, opts ...AppModuleOption) error {
	// 保存每个模块被依赖的次数
	inDegree := make(map[string]int)

	options := AppModuleOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	initModulesMap := make(map[string]string, 0)

	if len(options.modules) > 0 {
		for _, name := range options.modules {
			if _, ok := initModulesMap[name]; !ok {
				initModulesMap[name] = name
				module := g_appModules[name]
				for _, dep := range module.Depends {
					initModulesMap[dep] = dep
				}
			}
		}
	} else {
		for _, module := range g_appModules {
			initModulesMap[module.Name] = module.Name
			for _, dep := range module.Depends {
				initModulesMap[dep] = dep
			}
		}
	}

	// 构建依赖图
	for name := range initModulesMap {
		if _, ok := inDegree[name]; !ok {
			inDegree[name] = 0
		}
		module := g_appModules[name]
		if module == nil {
			return fmt.Errorf("module %s not found", name)
		}

		for _, dep := range module.Depends {
			if _, ok := inDegree[dep]; !ok {
				inDegree[dep] = 1
			} else {
				inDegree[dep]++
			}
		}
	}
	//logger.Debug("inDegree: ", inDegree)
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

		//logger.Debug("sortedModules: ", sortedModules)
	}

	// 检查是否有环
	if len(sortedModules) != len(initModulesMap) {
		logger.Debug("sortedModules: ", sortedModules)
		logger.Debug("appModules: ", initModulesMap)
		return errors.New("circular dependency detected")
	}

	// 按顺序初始化模块
	for _, name := range sortedModules {
		module := g_appModules[name]
		moduleConfig := appConfig
		if module.Options.config != nil {
			loadModuleConfig(module.Options.config)
		}

		err := module.Module.Init(moduleConfig, module.Depends)
		if err != nil {
			return fmt.Errorf("failed to initialize module %s: %w", name, err)
		}
		module.Ready = true
	}

	return nil
}
