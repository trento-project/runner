package runner

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

//go:embed ansible
var ansibleFS embed.FS

const (
	AnsibleMain       = "ansible/check.yml"
	AnsibleMeta       = "ansible/meta.yml"
	AnsibleConfigFile = "ansible/ansible.cfg"
	AnsibleHostFile   = "ansible/ansible_hosts"
)

//go:generate mockery --name=RunnerService --inpackage --filename=runner_mock.go

type RunnerService interface {
	IsCatalogReady() bool
	BuildCatalog() error
	GetCatalog() map[string]*Catalog
}

type runnerService struct {
	config  *Config
	catalog map[string]*Catalog
	ready   bool
}

func NewRunnerService(config *Config) (*runnerService, error) {
	runner := &runnerService{
		config: config,
		ready:  false,
	}

	return runner, nil
}

func (c *runnerService) IsCatalogReady() bool {
	return c.ready
}

func (c *runnerService) BuildCatalog() error {
	if err := createAnsibleFiles(c.config.AnsibleFolder); err != nil {
		return err
	}

	metaRunner, err := NewAnsibleMetaRunner(c.config)
	if err != nil {
		return err
	}

	// The checks catalog metadata playbook creates the checks catalog in the provider file path
	if err = metaRunner.RunPlaybook(); err != nil {
		log.Errorf("Error running the catalog meta-playbook")
		return err
	}

	// After the playbook is done, recover back the file content
	catalogRaw, err := ioutil.ReadFile(metaRunner.Envs[CatalogDestination])
	if err != nil {
		log.Fatal("Error when opening the catalog file: ", err)
	}

	var catalog map[string]*Catalog
	err = json.Unmarshal(catalogRaw, &catalog)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	c.catalog = catalog
	c.ready = true

	return nil
}

func (c *runnerService) GetCatalog() map[string]*Catalog {
	return c.catalog
}

func createAnsibleFiles(folder string) error {
	log.Infof("Creating the ansible file structure in %s", folder)
	// Clean the folder if it stores old files
	ansibleFolder := path.Join(folder, "ansible")
	err := os.RemoveAll(ansibleFolder)
	if err != nil {
		log.Error(err)
		return err
	}

	err = os.MkdirAll(ansibleFolder, 0755)
	if err != nil {
		log.Error(err)
		return err
	}

	// Create the ansible file structure from the FS
	err = fs.WalkDir(ansibleFS, "ansible", func(fileName string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dir.IsDir() {
			content, err := ansibleFS.ReadFile(fileName)
			if err != nil {
				log.Errorf("Error reading file %s", fileName)
				return err
			}
			f, err := os.Create(path.Join(folder, fileName))
			if err != nil {
				log.Errorf("Error creating file %s", fileName)
				return err
			}
			fmt.Fprintf(f, "%s", content)
		} else {
			os.Mkdir(path.Join(folder, fileName), 0755)
		}
		return nil
	})

	if err != nil {
		log.Errorf("An error ocurred during the ansible file structure creation: %s", err)
		return err
	}

	log.Info("Ansible file structure successfully created")

	return nil
}

func NewAnsibleMetaRunner(config *Config) (*AnsibleRunner, error) {
	playbookPath := path.Join(config.AnsibleFolder, AnsibleMeta)
	ansibleRunner := DefaultAnsibleRunner()

	if err := ansibleRunner.SetPlaybook(playbookPath); err != nil {
		return ansibleRunner, err
	}

	configFile := path.Join(config.AnsibleFolder, AnsibleConfigFile)
	ansibleRunner.SetConfigFile(configFile)
	destination := path.Join(config.AnsibleFolder, CatalogDestinationFile)
	ansibleRunner.SetCatalogDestination(destination)

	return ansibleRunner, nil
}

func NewAnsibleCheckRunner(config *Config) (*AnsibleRunner, error) {
	playbookPath := path.Join(config.AnsibleFolder, AnsibleMain)

	ansibleRunner := DefaultAnsibleRunner()

	if err := ansibleRunner.SetPlaybook(playbookPath); err != nil {
		return ansibleRunner, err
	}

	ansibleRunner.Check = true
	configFile := path.Join(config.AnsibleFolder, AnsibleConfigFile)
	ansibleRunner.SetConfigFile(configFile)
	ansibleRunner.SetTrentoApiData(config.ApiHost, config.ApiPort)

	return ansibleRunner, nil
}
