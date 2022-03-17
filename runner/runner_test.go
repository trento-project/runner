package runner

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	TestAnsibleFolder string = "../test/ansible_test"
)

type RunnerTestCase struct {
	suite.Suite
}

func TestRunnerTestCase(t *testing.T) {
	suite.Run(t, new(RunnerTestCase))
}

// TODO: This test could be improved to check the definitve ansible files structure
// once we have something fixed
func (suite *ApiTestCase) Test_CreateAnsibleFiles() {
	tmpDir, _ := ioutil.TempDir(os.TempDir(), "trentotest")
	err := createAnsibleFiles(tmpDir)

	suite.DirExists(path.Join(tmpDir, "ansible"))
	suite.NoError(err)

	os.RemoveAll(tmpDir)
}

func (suite *ApiTestCase) Test_NewAnsibleMetaRunner() {

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

func (suite *ApiTestCase) Test_NewAnsibleCheckRunner() {

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
