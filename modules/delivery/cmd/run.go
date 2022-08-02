package cmd

import (
	"fmt"
	"log"
	config "todo-clone/modules/confing"
	"todo-clone/modules/delivery/rest"
	"todo-clone/modules/domains"
	"todo-clone/modules/repository"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Version: config.AppVersion,
	Short:   fmt.Sprintf("%s Run", config.AppName),
	Run:     runModules,
}

func runModules(cmd *cobra.Command, args []string) {
	f := func(repo *repository.Repository) {
		if err := repo.AutoMigrate(&domains.Item{}); err != nil {
			log.Fatal(err)
		}
	}
	modules := fx.Options(config.Modules, rest.Modules, repository.Modules, fx.Invoke(f))
	fx.New(modules).Run()
}
