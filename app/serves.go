package app

import (
	"context"
	"fmt"
)

type IAppServe interface {
	Done() <-chan struct{}
	Start(context.Context) error
	Ready() bool
	Stop() error
}

type AppServe struct {
	Name   string
	Module string
	Serve  IAppServe
	//Done    chan Signal
}

var g_appServes = make(map[string]AppServe, 10)

func GetServeNames() []string {
	var names []string
	for k := range g_appServes {
		names = append(names, k)
	}
	return names
}

func GetServes() map[string]AppServe {
	return g_appServes
}

func GetServe(name string) (AppServe, bool) {
	v, ok := g_appServes[name]
	return v, ok
}

func GetServeNamesByModule(module string) []string {
	var names []string
	for _, v := range g_appServes {
		if v.Module == module {
			names = append(names, v.Name)
		}
	}
	return names
}

func RegisterServe(name string, module string, serve IAppServe) error {
	if _, ok := g_appModules[module]; !ok {
		return fmt.Errorf("module %s not found", module)
	}

	id := fmt.Sprintf("%s.%s", module, name)
	g_appServes[id] = AppServe{
		Name:   name,
		Module: module,
		Serve:  serve,
		//Done:    make(chan Signal),
	}

	return nil
}

func StartServe(ctx context.Context, module string, serve string) error {
	if v, ok := g_appModules[module]; !ok {
		return fmt.Errorf("module %s not found", module)
	} else {
		if !v.Ready {
			return fmt.Errorf("module %s not ready", module)
		}
	}
	serveId := fmt.Sprintf("%s.%s", module, serve)
	if v, ok := g_appServes[serveId]; ok {
		err := v.Serve.Start(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("serve %s not found", serveId)
}

func StartServes(ctx context.Context, serves []string) error {
	for _, v := range serves {
		if v, ok := g_appServes[v]; ok {
			err := v.Serve.Start(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func StopServes() error {
	for _, v := range g_appServes {
		if !v.Serve.Ready() {
			continue
		}

		err := v.Serve.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}
