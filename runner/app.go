package runner

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Host          string
	Port          int
	CallbacksUrl  string
	AnsibleFolder string
}

type App struct {
	config *Config
	Dependencies
}

type Dependencies struct {
	webEngine           *gin.Engine
	executionWorkerPool *ExecutionWorkerPool
	runnerService       RunnerService
}

func DefaultDependencies(config *Config) Dependencies {
	webEngine := gin.New()
	webEngine.Use(gin.Recovery())

	mode := os.Getenv(gin.EnvGinMode)
	gin.SetMode(mode)

	runnerService, err := NewRunnerService(config)
	if err != nil {
		log.Fatalf("Failed to create the runner instance: %s", err)
	}

	executionWorkerPool := NewExecutionWorkerPool(runnerService)

	return Dependencies{
		webEngine,
		executionWorkerPool,
		runnerService,
	}
}

func NewApp(config *Config) (*App, error) {
	return NewAppWithDeps(config, DefaultDependencies(config))
}

func NewAppWithDeps(config *Config, deps Dependencies) (*App, error) {
	app := &App{
		config:       config,
		Dependencies: deps,
	}

	apiGroup := deps.webEngine.Group("/api")
	{
		apiGroup.GET("/health", HealthHandler)
		apiGroup.GET("/ready", ReadyHandler(deps.runnerService))
		apiGroup.GET("/catalog", CatalogHandler(deps.runnerService))
		apiGroup.POST("/execute", ExecutionHandler(deps.runnerService))
	}

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", a.config.Host, a.config.Port)
	webServer := &http.Server{
		Addr:           address,
		Handler:        a.webEngine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	g, ctx := errgroup.WithContext(ctx)

	log.Infof("Starting web server at %s", address)
	g.Go(func() error {
		err := webServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	log.Infof("Starting execution requests worker pool....")
	g.Go(func() error {
		a.executionWorkerPool.Run(ctx)
		return nil
	})

	log.Infof("Building catalog....")
	g.Go(func() error {
		err := a.runnerService.BuildCatalog()
		if err != nil {
			return err
		}
		return nil
	})

	go func() {
		<-ctx.Done()
		log.Info("Web server is shutting down.")
		webServer.Close()
	}()

	return g.Wait()
}
