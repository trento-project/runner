package cmd

import (
	"github.com/spf13/viper"
	"github.com/trento-project/runner/runner"
)

func LoadConfig() *runner.Config {
	return &runner.Config{
		Host:          viper.GetString("host"),
		Port:          viper.GetInt("port"),
		CallbacksUrl:  viper.GetString("callbacks-url"),
		AnsibleFolder: viper.GetString("ansible-folder"),
	}
}
