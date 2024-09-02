package app

import (
	"github.com/spf13/cobra"
)

var g_Cmds = []*cobra.Command{}

func RegisterCmd(cmd *cobra.Command) {
	g_Cmds = append(g_Cmds, cmd)
}

func InitializeCmds(parent *cobra.Command) error {
	parent.AddCommand(
		g_Cmds...,
	)
	return nil
}
