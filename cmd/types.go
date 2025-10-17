package cmd

type CmdOptions struct {
	AppName     string `json:"app_name"`
	AppVersion  string `json:"app_version"`
	EnvPrefix   string `json:"env_prefix"`
	Description string `json:"description"`
}

func NewCmdOptions() *CmdOptions {
	return &CmdOptions{
		AppName:    "app",
		AppVersion: "0.0.1",
		EnvPrefix:  "APP",
	}
}

func (o *CmdOptions) WithAppName(appName string) *CmdOptions {
	o.AppName = appName
	return o
}
func (o *CmdOptions) WithAppVersion(appVersion string) *CmdOptions {
	o.AppVersion = appVersion
	return o
}
func (o *CmdOptions) WithEnvPrefix(envPrefix string) *CmdOptions {
	o.EnvPrefix = envPrefix
	return o
}

func (o *CmdOptions) WithDescription(description string) *CmdOptions {
	o.Description = description
	return o
}
