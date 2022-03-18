package runner

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HealthApiTestCase struct {
	suite.Suite
	config *Config
}

func TestHealthApiTestCase(t *testing.T) {
	suite.Run(t, new(HealthApiTestCase))
}

func (suite *HealthApiTestCase) SetupTest() {
	suite.config = &Config{}
}

func (suite *HealthApiTestCase) Test_ApiHealthTest() {
	app, err := NewApp(suite.config)
	if err != nil {
		suite.T().Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/health", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedJson, _ := json.Marshal(map[string]string{"status": "ok"})
	suite.Equal(200, resp.Code)
	suite.JSONEq(string(expectedJson), resp.Body.String())
}

func (suite *HealthApiTestCase) Test_ApiReadyTest() {
	mockRunnerService := new(MockRunnerService)
	mockRunnerService.On("IsCatalogReady").Return(true)

	deps := Dependencies{
		mockRunnerService,
	}

	app, err := NewAppWithDeps(suite.config, deps)
	if err != nil {
		suite.T().Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/ready", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedJson, _ := json.Marshal(map[string]bool{"ready": true})
	suite.Equal(200, resp.Code)
	suite.JSONEq(string(expectedJson), resp.Body.String())
}
