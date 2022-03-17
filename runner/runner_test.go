package runner

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/trento-project/runner/runner/mocks"
)

const (
	TestAnsibleFolder string = "../test/ansible_test"
)

type RunnerTestCase struct {
	suite.Suite
	runnerService RunnerService
	ansibleDir    string
}

func TestRunnerTestCase(t *testing.T) {
	suite.Run(t, new(RunnerTestCase))
}

func (suite *RunnerTestCase) SetupTest() {
	tmpDir, _ := ioutil.TempDir(os.TempDir(), "trentotest")
	runnerService, _ := NewRunnerService(&Config{AnsibleFolder: tmpDir})
	suite.runnerService = runnerService
	suite.ansibleDir = tmpDir
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
		ApiHost:       "127.0.0.1",
		ApiPort:       8000,
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
		ApiHost:       "127.0.0.1",
		ApiPort:       8000,
		AnsibleFolder: TestAnsibleFolder,
	}

	a, err := NewAnsibleCheckRunner(cfg)

	expectedMetaRunner := &AnsibleRunner{
		Playbook: path.Join(TestAnsibleFolder, "ansible/check.yml"),
		Envs: map[string]string{
			"ANSIBLE_CONFIG":      path.Join(TestAnsibleFolder, "ansible/ansible.cfg"),
			"TRENTO_WEB_API_HOST": "127.0.0.1",
			"TRENTO_WEB_API_PORT": "8000",
		},
		Check: true,
	}

	suite.NoError(err)
	suite.Equal(expectedMetaRunner, a)
}
