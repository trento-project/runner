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

func (suite *CatalogApiTestCase) Test_GetCatalogTest_NotReady() {
	mockRunnerService := new(MockRunnerService)
	mockRunnerService.On("IsCatalogReady").Return(false)

	deps := setupTestDependencies()
	deps.runnerService = mockRunnerService

	app, err := NewAppWithDeps(suite.config, deps)
	if err != nil {
		suite.T().Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/catalog", nil)
	app.webEngine.ServeHTTP(resp, req)

	suite.Equal(204, resp.Code)
	suite.Equal("", resp.Body.String())
}

func (suite *CatalogApiTestCase) Test_GetCatalogTest_Ready() {
	returnedCatalog := &Catalog{
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
	}

	mockRunnerService := new(MockRunnerService)
	mockRunnerService.On("IsCatalogReady").Return(true)
	mockRunnerService.On("GetCatalog").Return(returnedCatalog)

	deps := setupTestDependencies()
	deps.runnerService = mockRunnerService

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
