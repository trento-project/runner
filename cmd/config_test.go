package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/runner/runner"
)

type RunnerCmdTestSuite struct {
	suite.Suite
	cmd *cobra.Command
}

func TestRunnerCmdTestSuite(t *testing.T) {
	suite.Run(t, new(RunnerCmdTestSuite))
}

func (suite *RunnerCmdTestSuite) SetupTest() {
	os.Clearenv()

	cmd := NewRunnerCmd()
	cmd.Run = func(cmd *cobra.Command, args []string) {
		// do nothing
	}

	cmd.Commands()[0].Run = func(cmd *cobra.Command, args []string) {
		// do nothing
	}

	cmd.SetArgs([]string{
		"start",
	})

	var b bytes.Buffer
	cmd.SetOut(&b)

	suite.cmd = cmd
}

func (suite *RunnerCmdTestSuite) TearDownTest() {
	suite.cmd.Execute()

	expectedConfig := &runner.Config{
		Host:          "localhost",
		Port:          5678,
		CallbacksUrl:  "http://192.168.1.1:8000/api/runner/callbacks",
		AnsibleFolder: "path/to/ansible",
	}
	config := LoadConfig()

	suite.EqualValues(expectedConfig, config)
}

func (suite *RunnerCmdTestSuite) TestConfigFromFlags() {
	suite.cmd.SetArgs([]string{
		"start",
		"--host=localhost",
		"--port=5678",
		"--callbacks-url=http://192.168.1.1:8000/api/runner/callbacks",
		"--ansible-folder=path/to/ansible",
	})
}

func (suite *RunnerCmdTestSuite) TestConfigFromEnv() {
	os.Setenv("TRENTO_RUNNER_HOST", "localhost")
	os.Setenv("TRENTO_RUNNER_PORT", "5678")
	os.Setenv("TRENTO_RUNNER_CALLBACKS_URL", "http://192.168.1.1:8000/api/runner/callbacks")
	os.Setenv("TRENTO_RUNNER_ANSIBLE_FOLDER", "path/to/ansible")
}

func (suite *RunnerCmdTestSuite) TestConfigFromFile() {
	os.Setenv("TRENTO_RUNNER_CONFIG", "../test/fixtures/config/runner.yaml")
}
