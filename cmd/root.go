package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/shin1x1/docker-base-chaser/internal/config"
	"github.com/shin1x1/docker-base-chaser/internal/handler"
	"github.com/spf13/cobra"
)

const CommandName = "docker-base-chaser"

var (
	configFilePath string
	dryRun         bool
)

func init() {
	rootCmd.Flags().StringVarP(&configFilePath, "config", "c", "base-chaser.yaml", "config file path(default: base-chaser.yaml")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "dry-run(default: false)")
}

var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [command] [command args]", CommandName),
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		var targets *handler.Targets
		l, err := config.LoadLock(config.CreateLockPath(configFilePath))
		if err == nil {
			targets = l.CreateTargets()
		} else {
			c, err := config.LoadConfig(configFilePath)
			if err != nil {
				return err
			}

			targets = c.CreateTargets(c)
		}

		h := handler.NewRootHandler(dryRun)
		if err := h.Exec(targets); err != nil {
			return err
		}

		// If dryRun is enabled, will not update lock file.
		if dryRun {
			return nil
		}

		now := time.Now()
		lock := config.NewLockWithTargets(targets, &now)

		return lock.Save("base-chaser.lock")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
