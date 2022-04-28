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

	expectedCatalog := &Catalog{
		&CatalogCheck{
			ID:             "156F64",
			Name:           "1.1.1",
			Group:          "Corosync",
			Provider:       "azure",
			Description:    "description azure",
			Remediation:    "remediation",
			Implementation: "implementation",
			Labels:         "generic",
			Premium:        false,
		},
		&CatalogCheck{
			ID:             "156F64",
			Name:           "1.1.1",
			Group:          "Corosync",
			Provider:       "dev",
			Description:    "description dev",
			Remediation:    "remediation",
			Implementation: "implementation",
			Labels:         "generic",
			Premium:        false,
		},
	}

	suite.NoError(err)
	suite.Equal(true, suite.runnerService.IsCatalogReady())
	suite.Equal(expectedCatalog, suite.runnerService.GetCatalog())
}

func (suite *RunnerTestCase) Test_ScheduleExecution() {
	execution := &ExecutionEvent{ExecutionID: uuid.New()}
	err := suite.runnerService.ScheduleExecution(execution)
	suite.NoError(err)
	suite.Equal(execution, <-suite.runnerService.GetChannel())
}

func (suite *RunnerTestCase) Test_ScheduleExecution_Full() {
	ch := suite.runnerService.GetChannel()
	for range [executionChannelSize]int64{} {
		ch <- &ExecutionEvent{ExecutionID: uuid.New()}
	}

	execution := &ExecutionEvent{ExecutionID: uuid.New()}
	err := suite.runnerService.ScheduleExecution(execution)
	suite.EqualError(err, "Cannot process more executions")
}

func (suite *RunnerTestCase) Test_Execute() {
	os.MkdirAll(path.Join(suite.ansibleDir, "ansible"), 0755)
	os.Create(path.Join(suite.ansibleDir, "ansible/check.yml"))
	defer os.RemoveAll(suite.ansibleDir)

	dummyID := uuid.New()
	clusterDummyID := uuid.New()
	executionStartedPayload := map[string]string{"cluster_id": clusterDummyID.String()}
	suite.callbacksClient.On(
		"Callback", dummyID, "execution_started", executionStartedPayload).Return(nil)

	cmd := exec.Command("ls") // Dummy command to execute something

	mockCommand := new(mocks.CustomCommand)
	customExecCommand = mockCommand.Execute

	mockCommand.On(
		"Execute",
		"ansible-playbook",
		path.Join(suite.ansibleDir, "ansible/check.yml"),
		fmt.Sprintf("--inventory=%s/ansible/inventories/%s/ansible_hosts", suite.ansibleDir, dummyID.String()),
		"--check",
	).Return(cmd)

	execution := &ExecutionEvent{ExecutionID: dummyID, ClusterID: clusterDummyID}
	err := suite.runnerService.Execute(execution)

	suite.NoError(err)
}

func (suite *RunnerTestCase) Test_Execute_CallbackError() {
	dummyID := uuid.New()
	clusterDummyID := uuid.New()
	expectedError := fmt.Errorf("error running callback")
	executionStartedPayload := map[string]string{"cluster_id": clusterDummyID.String()}
	suite.callbacksClient.On(
		"Callback", dummyID, "execution_started", executionStartedPayload).Return(expectedError)

	execution := &ExecutionEvent{ExecutionID: dummyID, ClusterID: clusterDummyID}
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
	tmpDir, _ := ioutil.TempDir(os.TempDir(), "trentotest")
	os.MkdirAll(path.Join(tmpDir, "ansible"), 0755)
	os.Create(path.Join(tmpDir, "ansible/check.yml"))
	defer os.RemoveAll(tmpDir)

	cfg := &Config{
		CallbacksUrl:  "http://192.168.1.1:8000/api/runner/callbacks",
		AnsibleFolder: tmpDir,
	}

	executionID := uuid.New()
	clusterID := uuid.New()
	host1ID := uuid.New()
	host2ID := uuid.New()
	executionEvent := &ExecutionEvent{
		ExecutionID: executionID,
		ClusterID:   clusterID,
		Provider:    "azure",
		Checks:      []string{"check1", "check2"},
		Hosts: []*Host{
			&Host{
				HostID:  host1ID,
				Address: "192.168.10.1",
				User:    "user1",
			},
			&Host{
				HostID:  host2ID,
				Address: "192.168.10.2",
				User:    "user2",
			},
		},
	}

	a, err := NewAnsibleCheckRunner(cfg, executionEvent)

	inventoryFile := path.Join(tmpDir, fmt.Sprintf("ansible/inventories/%s/ansible_hosts", executionID.String()))

	expectedChecksRunner := &AnsibleRunner{
		Playbook:  path.Join(tmpDir, "ansible/check.yml"),
		Inventory: inventoryFile,
		Envs: map[string]string{
			"ANSIBLE_CONFIG":       path.Join(tmpDir, "ansible/ansible.cfg"),
			"TRENTO_CALLBACKS_URL": "http://192.168.1.1:8000/api/runner/callbacks",
			"TRENTO_EXECUTION_ID":  executionID.String(),
		},
		Check: true,
	}

	inventoryContent, err := ioutil.ReadFile(inventoryFile)
	expectedFile := "\n" +
		"[%s]\n" +
		"%s ansible_host=192.168.10.1 ansible_user=user1 cluster_selected_checks='[\"check1\",\"check2\"]' provider=azure \n" +
		"%s ansible_host=192.168.10.2 ansible_user=user2 cluster_selected_checks='[\"check1\",\"check2\"]' provider=azure \n"

	suite.NoError(err)
	suite.Equal(expectedChecksRunner, a)
	suite.Equal(fmt.Sprintf(expectedFile, clusterID.String(), host1ID.String(), host2ID.String()), string(inventoryContent))
}
