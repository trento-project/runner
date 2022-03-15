package cmd

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/trento-project/runner/internal"
	"github.com/trento-project/runner/runner"
)

func NewRunnerCmd() *cobra.Command {
	var cfgFile string
	var logLevel string

	runnerCmd := &cobra.Command{
		Use:   "trento-runner",
		Short: "Trento Runner component, which is responsible of running the Trento checks",
		PersistentPreRunE: func(runnerCmd *cobra.Command, _ []string) error {
			return internal.InitConfig("runner")
		},
	}

	runnerCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.trento/runner.yaml)")
	runnerCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "then minimum severity (error, warn, info, debug) of logs to output")

	runnerCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})

	addStartCmd(runnerCmd)
	addVersionCmd(runnerCmd)

	return runnerCmd
}

func Execute() {
	if err := NewRunnerCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func addStartCmd(runnerCmd *cobra.Command) {
	var apiHost string
	var apiPort int
	var interval int
	var ansibleFolder string

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Starts the runner process. This process takes care of running the checks",
		Run:   start,
	}

	startCmd.Flags().StringVar(&apiHost, "api-host", "0.0.0.0", "Trento web server API host")
	startCmd.Flags().IntVar(&apiPort, "api-port", 8080, "Trento web server API port")
	startCmd.Flags().IntVarP(&interval, "interval", "i", 5, "Interval in minutes to run the checks")
	startCmd.Flags().StringVar(&ansibleFolder, "ansible-folder", "/tmp/trento", "Folder where the ansible file structure will be created")

	runnerCmd.AddCommand(startCmd)
}

func start(*cobra.Command, []string) {
	var err error

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	config := LoadConfig()

	runner, err := runner.NewRunner(config)
	if err != nil {
		log.Fatal("Failed to create the runner instance: ", err)
	}

	go func() {
		quit := <-signals
		log.Printf("Caught %s signal!", quit)

		log.Println("Stopping the runner...")
		runner.Stop()
	}()

	err = runner.Start()
	if err != nil {
		log.Fatal("Failed to start the runner service: ", err)
	}
}
