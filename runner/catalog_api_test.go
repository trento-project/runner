package runner

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CatalogApiTestCase struct {
	suite.Suite
	config *Config
}

func TestCatalogApiTestCase(t *testing.T) {
	suite.Run(t, new(CatalogApiTestCase))
}

func (suite *CatalogApiTestCase) SetupTest() {
	suite.config = &Config{}
}

func (suite *CatalogApiTestCase) Test_GetCatalogTest() {
	returnedCatalog := map[string]*Catalog{
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
	}

	mockRunnerService := new(MockRunnerService)
	mockRunnerService.On("GetCatalog").Return(returnedCatalog)

	deps := Dependencies{
		mockRunnerService,
	}

	app, err := NewAppWithDeps(suite.config, deps)
	if err != nil {
		suite.T().Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/catalog", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedJson, _ := json.Marshal(returnedCatalog)
	suite.Equal(200, resp.Code)
	suite.JSONEq(string(expectedJson), resp.Body.String())
}
