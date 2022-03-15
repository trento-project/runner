package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ClusterSettings struct {
	ID             string            `json:"id"`
	SelectedChecks []string          `json:"selected_checks"`
	Hosts          []*HostConnection `json:"hosts"`
}

type HostConnection struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	User    string `json:"user"`
}

type ClustersSettingsResponse []*ClusterSettings

func (t *trentoApiService) GetClustersSettings() (ClustersSettingsResponse, error) {
	body, statusCode, err := t.getJson("clusters/settings")
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error during the request with status code %d", statusCode)
	}

	var clustersSettings ClustersSettingsResponse

	err = json.Unmarshal(body, &clustersSettings)
	if err != nil {
		return nil, err
	}

	return clustersSettings, nil
}
