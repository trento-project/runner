package cmd

import (
	"context"
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
			runnerCmd.Flags().VisitAll(func(f *pflag.Flag) {
				viper.BindPFlag(f.Name, f)
			})

			runnerCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
				viper.BindPFlag(f.Name, f)
			})

			return internal.InitConfig("runner")
		},
	}

	runnerCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.trento/runner.yaml)")
	runnerCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "then minimum severity (error, warn, info, debug) of logs to output")

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
	var host string
	var port int
	var callbacksUrl string
	var ansibleFolder string

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Starts the runner process. This process takes care of running the checks",
		Run:   start,
	}

	startCmd.Flags().StringVar(&host, "host", "0.0.0.0", "Trento Runner API host")
	startCmd.Flags().IntVar(&port, "port", 8080, "Trento Runner API port")
	startCmd.Flags().StringVar(&callbacksUrl, "callbacks-url", "", "Trento web server runner callbacks API url")
	startCmd.MarkFlagRequired("callbacks-url")
	startCmd.Flags().StringVar(&ansibleFolder, "ansible-folder", "/tmp/trento", "Folder where the ansible file structure will be created")

	runnerCmd.AddCommand(startCmd)
}

func start(*cobra.Command, []string) {
	var err error

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	config := LoadConfig()

	app, err := runner.NewApp(config)
	if err != nil {
		log.Fatal("Failed to create the runner application: ", err)
	}

	go func() {
		quit := <-signals
		log.Printf("Caught %s signal!", quit)

		log.Println("Stopping the runner...")
		cancel()
	}()

	if err = app.Start(ctx); err != nil {
		log.Fatal("Failed to start the runner application: ", err)
	}
}
