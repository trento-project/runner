package cmd

import (
	"time"

	"github.com/spf13/viper"
	"github.com/trento-project/runner/runner"
)

func LoadConfig() *runner.Config {
	return &runner.Config{
		ApiHost:       viper.GetString("api-host"),
		ApiPort:       viper.GetInt("api-port"),
		Interval:      time.Duration(viper.GetInt("interval")) * time.Minute,
		AnsibleFolder: viper.GetString("ansible-folder"),
	}
}
