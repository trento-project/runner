package runner

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type InventoryTestSuite struct {
	suite.Suite
}

func TestInventoryTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryTestSuite))
}

func (suite *InventoryTestSuite) Test_CreateInventory() {
	tmpDir, _ := ioutil.TempDir(os.TempDir(), "trentotest")
	destination := path.Join(tmpDir, "ansible_hosts")

	content := &InventoryContent{
		Nodes: []*Node{
			&Node{
				Name:        "node1",
				AnsibleHost: "192.168.10.1",
				AnsibleUser: "trento",
				Variables: map[string]interface{}{
					"key1": "value1",
					"key2": []string{"value2", "value3"},
				},
			},
			&Node{
				Name: "node2",
			},
		},
		Groups: []*Group{
			&Group{
				Name: "group1",
				Nodes: []*Node{
					{
						Name:        "node3",
						AnsibleHost: "192.168.11.1",
						AnsibleUser: "trento",
						Variables: map[string]interface{}{
							"key1": 1,
							"key2": []string{"value2", "value3"},
						},
					},
					&Node{
						Name: "node4",
					},
				},
			},
			&Group{
				Name: "group2",
				Nodes: []*Node{
					{
						Name: "node5",
					},
					&Node{
						Name: "node6",
					},
				},
			},
		},
	}

	err := CreateInventory(destination, content)

	suite.NoError(err)
	suite.FileExists(destination)

	// Cannot use backticks as the lines have a final space in many lines
	expectedContent := "\n" +
		"node1 ansible_host=192.168.10.1 ansible_user=trento key1=value1 key2=[value2 value3] \n" +
		"node2 ansible_host= ansible_user= \n" +
		"[group1]\n" +
		"node3 ansible_host=192.168.11.1 ansible_user=trento key1=1 key2=[value2 value3] \n" +
		"node4 ansible_host= ansible_user= \n" +
		"[group2]\n" +
		"node5 ansible_host= ansible_user= \n" +
		"node6 ansible_host= ansible_user= \n"

	data, err := ioutil.ReadFile(destination)
	if err == nil {
		suite.Equal(expectedContent, string(data))
	}
}

func (suite *InventoryTestSuite) Test_NewClusterInventoryContent() {
	cluster := uuid.New()
	host1 := uuid.New()
	host2 := uuid.New()

	executionEvent := &ExecutionEvent{
		ExecutionID: uuid.New(),
		ClusterID:   cluster,
		Provider:    "azure",
		Checks:      []string{"check1", "check2"},
		Hosts: []*Host{
			&Host{
				HostID:  host1,
				Address: "192.168.10.1",
				User:    "user1",
			},
			&Host{
				HostID:  host2,
				Address: "192.168.10.2",
				User:    "user2",
			},
		},
	}

	content, err := NewClusterInventoryContent(executionEvent)

	expectedContent := &InventoryContent{
		Groups: []*Group{
			&Group{
				Name: cluster.String(),
				Nodes: []*Node{
					&Node{
						Name: host1.String(),
						Variables: map[string]interface{}{
							"cluster_selected_checks": "'[\"check1\",\"check2\"]'",
							"provider":                "azure",
						},
						AnsibleHost: "192.168.10.1",
						AnsibleUser: "user1",
					},
					&Node{
						Name: host2.String(),
						Variables: map[string]interface{}{
							"cluster_selected_checks": "'[\"check1\",\"check2\"]'",
							"provider":                "azure",
						},
						AnsibleHost: "192.168.10.2",
						AnsibleUser: "user2",
					},
				},
			},
		},
	}

	suite.NoError(err)
	suite.ElementsMatch(expectedContent.Groups, content.Groups)
}
