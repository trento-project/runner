package runner

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/trento-project/runner/runner/mocks"
)

const (
	TestAnsibleFolder string = "../test/ansible_test"
)

type RunnerTestCase struct {
	suite.Suite
	runnerService   RunnerService
	ansibleDir      string
	callbacksClient *mocks.CallbacksClient
}

func TestRunnerTestCase(t *testing.T) {
	suite.Run(t, new(RunnerTestCase))
}

func (suite *RunnerTestCase) SetupTest() {
	callbacksClient := new(mocks.CallbacksClient)
	tmpDir, _ := ioutil.TempDir(os.TempDir(), "trentotest")
	runnerService, _ := NewRunnerService(&Config{AnsibleFolder: tmpDir})
	runnerService.callbacksClient = callbacksClient
	suite.runnerService = runnerService
	suite.ansibleDir = tmpDir
	suite.callbacksClient = callbacksClient
}

func (suite *RunnerTestCase) Test_BuildCatalog() {
	suite.Equal(false, suite.runnerService.IsCatalogReady())

	cmd := exec.Command("cp", "../test/fixtures/catalog.json", path.Join(suite.ansibleDir, "ansible"))

	mockCommand := new(mocks.CustomCommand)
	customExecCommand = mockCommand.Execute

	mockCommand.On(
		"Execute", "ansible-playbook", path.Join(suite.ansibleDir, "ansible/meta.yml")).Return(
		cmd,
	)

	err := suite.runnerService.BuildCatalog()

	expectedMap := map[string]*Catalog{
		"azure": &Catalog{
			Checks: []*CatalogCheck{
				{
					ID:             "156F64",
					Name:           "1.1.1",
					Group:          "Corosync",
					Description:    "description azure",
					Remediation:    "remediation",
					Implementation: "implementation",
					Labels:         "generic",
					Premium:        false,
				},
			},
		},
		"dev": &Catalog{
			Checks: []*CatalogCheck{
				{
					ID:             "156F64",
					Name:           "1.1.1",
					Group:          "Corosync",
					Description:    "description dev",
					Remediation:    "remediation",
					Implementation: "implementation",
					Labels:         "generic",
					Premium:        false,
				},
			},
		},
	}

	suite.NoError(err)
	suite.Equal(true, suite.runnerService.IsCatalogReady())
	suite.Equal(expectedMap, suite.runnerService.GetCatalog())
}

func (suite *RunnerTestCase) Test_ScheduleExecution() {
	execution := &ExecutionEvent{ID: uuid.New()}
	err := suite.runnerService.ScheduleExecution(execution)
	suite.NoError(err)
	suite.Equal(execution, <-suite.runnerService.GetChannel())
}

func (suite *RunnerTestCase) Test_ScheduleExecution_Full() {
	ch := suite.runnerService.GetChannel()
	for range [executionChannelSize]int64{} {
		ch <- &ExecutionEvent{ID: uuid.New()}
	}

	execution := &ExecutionEvent{ID: uuid.New()}
	err := suite.runnerService.ScheduleExecution(execution)
	suite.EqualError(err, "Cannot process more executions")
}

func (suite *RunnerTestCase) Test_Execute() {
	dummyID := uuid.New()
	suite.callbacksClient.On("Callback", dummyID, "execution_started", nil).Return(nil)

	execution := &ExecutionEvent{ID: dummyID}
	err := suite.runnerService.Execute(execution)

	suite.NoError(err)
}

func (suite *RunnerTestCase) Test_Execute_CallbackError() {
	dummyID := uuid.New()
	expectedError := fmt.Errorf("error running callback")
	suite.callbacksClient.On("Callback", dummyID, "execution_started", nil).Return(expectedError)

	execution := &ExecutionEvent{ID: dummyID}
	err := suite.runnerService.Execute(execution)

	suite.EqualError(err, expectedError.Error())
}

// TODO: This test could be improved to check the definitve ansible files structure
// once we have something fixed
func (suite *RunnerTestCase) Test_CreateAnsibleFiles() {
	tmpDir, _ := ioutil.TempDir(os.TempDir(), "trentotest")
	err := createAnsibleFiles(tmpDir)

	suite.DirExists(path.Join(tmpDir, "ansible"))
	suite.NoError(err)

	os.RemoveAll(tmpDir)
}

func (suite *RunnerTestCase) Test_NewAnsibleMetaRunner() {

	cfg := &Config{
		CallbacksUrl:  "http://192.168.1.1:8000/api/runner/callbacks",
		AnsibleFolder: TestAnsibleFolder,
	}

	a, err := NewAnsibleMetaRunner(cfg)

	expectedMetaRunner := &AnsibleRunner{
		Playbook: path.Join(TestAnsibleFolder, "ansible/meta.yml"),
		Envs: map[string]string{
			"ANSIBLE_CONFIG":      path.Join(TestAnsibleFolder, "ansible/ansible.cfg"),
			"CATALOG_DESTINATION": path.Join(TestAnsibleFolder, "ansible/catalog.json"),
		},
		Check: false,
	}

	suite.NoError(err)
	suite.Equal(expectedMetaRunner, a)
}

func (suite *RunnerTestCase) Test_NewAnsibleCheckRunner() {

	cfg := &Config{
		CallbacksUrl:  "http://192.168.1.1:8000/api/runner/callbacks",
		AnsibleFolder: TestAnsibleFolder,
	}

	a, err := NewAnsibleCheckRunner(cfg)

	expectedMetaRunner := &AnsibleRunner{
		Playbook: path.Join(TestAnsibleFolder, "ansible/check.yml"),
		Envs: map[string]string{
			"ANSIBLE_CONFIG":       path.Join(TestAnsibleFolder, "ansible/ansible.cfg"),
			"TRENTO_CALLBACKS_URL": "http://192.168.1.1:8000/api/runner/callbacks",
		},
		Check: true,
	}

	suite.NoError(err)
	suite.Equal(expectedMetaRunner, a)
}
