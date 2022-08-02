package cmd

import (
	"fmt"
	"log"
	config "todo-clone/modules/confing"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Version: config.AppVersion,
	Short:   fmt.Sprintf("%s Config", config.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		log.Print(config.JSON()) // TODO: config.JSON()
	},
}
