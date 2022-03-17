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
	ApiHost       string
	ApiPort       int
	Interval      time.Duration
	AnsibleFolder string
}

type App struct {
	config    *Config
	webEngine *gin.Engine
	deps      Dependencies
}

type Dependencies struct {
	runnerService RunnerService
}

func DefaultDependencies(config *Config) Dependencies {
	runnerService, err := NewRunnerService(config)
	if err != nil {
		log.Fatalf("Failed to create the runner instance: %s", err)
	}

	return Dependencies{
		runnerService: runnerService,
	}
}

func NewApp(config *Config) (*App, error) {
	return NewAppWithDeps(config, DefaultDependencies(config))
}

func NewAppWithDeps(config *Config, deps Dependencies) (*App, error) {
	app := &App{
		config: config,
		deps:   deps,
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	mode := os.Getenv(gin.EnvGinMode)
	gin.SetMode(mode)

	apiGroup := engine.Group("/api")
	{
		apiGroup.GET("/health", HealthHandler)
		apiGroup.GET("/ready", ReadyHandler(deps.runnerService))
		apiGroup.GET("/catalog", CatalogHandler(deps.runnerService))
	}

	app.webEngine = engine

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

	log.Infof("Building catalog....")
	g.Go(func() error {
		err := a.deps.runnerService.BuildCatalog()
		if err != nil {
			return err
		}
		return nil
	})

	/*
		log.Infof("Starting runner....")
		g.Go(func() error {
			err := a.deps.runnerService.Start(ctx)
			if err != nil {
				return err
			}
			return nil
		})
	*/

	go func() {
		<-ctx.Done()
		log.Info("Web server is shutting down.")
		webServer.Close()
	}()

	return g.Wait()
}
