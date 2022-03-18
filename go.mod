module github.com/trento-project/runner

go 1.16

replace github.com/trento-project/runner => ./

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.1
	github.com/swaggo/swag v1.8.0
	github.com/vektra/mockery/v2 v2.10.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)