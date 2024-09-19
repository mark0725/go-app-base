package app

import (
	"github.com/spf13/cobra"
)

var g_Cmds map[string]*AppCmd = make(map[string]*AppCmd)

type AppCmdOption func(*AppCmdOptions)

type AppCmd struct {
	Name    string
	Modules []string
	Cmd     *cobra.Command
}

type AppCmdOptions struct {
}

func RegisterCmd(name string, cmd *cobra.Command, modules []string) {
	g_Cmds[name] = &AppCmd{
		Name:    name,
		Modules: modules,
		Cmd:     cmd,
	}
}

func InitializeCmds(cmd string, parent *cobra.Command, opts ...AppCmdOption) error {
	for _, v := range g_Cmds {
		parent.AddCommand(v.Cmd)
	}
	return nil
}
