package app

type IAppServe interface {
	Start() error
	Wait()
	Ready() bool
	Stop() error
}

type AppServe struct {
	Name  string
	Serve IAppServe
	//Done    chan Signal
}

var g_appServes = make(map[string]AppServe, 10)

func RegisterServe(name string, serve IAppServe) {
	g_appServes[name] = AppServe{
		Name:  name,
		Serve: serve,
		//Done:    make(chan Signal),
	}
}

func StartServes() error {
	for _, v := range g_appServes {
		err := v.Serve.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func StopServes() error {
	for _, v := range g_appServes {
		err := v.Serve.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}
