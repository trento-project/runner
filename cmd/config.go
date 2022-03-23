package cmd

import (
	"time"

	"github.com/spf13/viper"
	"github.com/trento-project/runner/runner"
)

func LoadConfig() *runner.Config {
	return &runner.Config{
		Host:          viper.GetString("host"),
		Port:          viper.GetInt("port"),
		CallbacksUrl:  viper.GetString("callbacks-url"),
		Interval:      time.Duration(viper.GetInt("interval")) * time.Minute,
		AnsibleFolder: viper.GetString("ansible-folder"),
	}
}
