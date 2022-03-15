package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/trento-project/runner/version"
)

func addVersionCmd(runnerCmd *cobra.Command) {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Trento Runner",
		Long:  `All software has versions. This is Trento Runner's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Trento Runner %s version %s\nbuilt with %s %s/%s\n", version.Flavor, version.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		},
	}

	runnerCmd.AddCommand(versionCmd)
}
