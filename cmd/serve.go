package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/mark0725/go-app-base/config"

	base_app "github.com/mark0725/go-app-base/app"
	"github.com/spf13/cobra"
)

func NewServeCommand(appConfig config.IAppConfig, options *CmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"start"},
		Short:   "Start " + options.AppName,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			app := base_app.NewApplication(appConfig, base_app.NewApplicationOptions().AppName(options.AppName).EnvPrefix(options.EnvPrefix).Version(options.AppVersion))
			err := app.AppInit()
			if err != nil {
				fmt.Print("App init error:", err)
				return err
			}

			ctx, cancel := context.WithCancel(context.Background())

			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

			done := make(chan bool, 1)
			go func() {
				if err := app.Run(ctx); err != nil {
					fmt.Printf("Server error: %v\n", err)
				}
				done <- true
			}()

			fmt.Println("Server is running. Press Ctrl+C to stop.")

			sig := <-signals
			fmt.Printf("Received signal: %s. Initiating graceful shutdown...\n", sig)

			cancel()

			select {
			case <-done:
				fmt.Println("Cleanup done. Shutdown completed.")
			case <-signals:
				fmt.Println("Forced shutdown.")
			case <-time.After(5 * time.Second):
				fmt.Println("shutdown completed.")
			}

			fmt.Println("Server stopped.")
			return nil
		},
	}
	return cmd
}
