package cmd

import (
	"fmt"
	"os"

	"github.com/mfojtik/oinc/pkg/log"
	"github.com/mfojtik/oinc/pkg/steps"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the OpenShift server in a container.",
	Long:  `Runs the OpenShift server in a container`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLogging()
		dirs := &steps.PrepareDirsStep{}
		if err := dirs.Execute(); err != nil {
			log.Critical("%s", err)
		}

		server := &steps.RunOpenShiftStep{}
		if err := server.Execute(); err != nil {
			log.Critical("%s", err)
		}

		(&steps.FixPermissionsStep{}).Execute()

		registry := &steps.InstallRegistryStep{}
		if err := registry.Execute(); err != nil {
			log.Critical("%s", err)
		}

		router := &steps.InstallRouterStep{}
		if err := router.Execute(); err != nil {
			log.Critical("%s", err)
		}

		templates := &steps.InstallTemplatesStep{}
		if err := templates.Execute(); err != nil {
			log.Critical("%s", err)
		}

		user := &steps.CreateUserStep{}
		if err := user.Execute(); err != nil {
			log.Critical("%s", err)
		}

		fmt.Fprintf(os.Stdout, `

OpenShift is now running! To access it using CLI tools, please run this command:

$ eval $(oinc env)
$ oc status

`)

	},
}

func init() {
	addPersistentFlags(RunCmd)
}
