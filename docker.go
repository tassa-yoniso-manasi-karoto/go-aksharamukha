
package aksharamukha

import (
	"time"
	"sync"

	"github.com/rs/zerolog"
	"github.com/gookit/color"
	"github.com/k0kubun/pp"
	
	"github.com/tassa-yoniso-manasi-karoto/dockerutil"
)

const (
	remote = "https://github.com/virtualvinodh/aksharamukha.git"
	projectName = "aksharamukha"
	// Main containers from the compose file
	containerFront = "aksharamukha-front-1"
	containerBack = "aksharamukha-back-1"
	containerFonts = "aksharamukha-fonts-1"
)

var (
	QueryTO = 1 * time.Hour
	dockerInstance *dockerutil.DockerManager
	dockerOnce sync.Once
	dockerMu sync.Mutex
)

type Aksharamukha struct {
	docker *dockerutil.DockerManager
	logger *dockerutil.ContainerLogConsumer
}

// NewDocker creates or returns an existing Docker manager instance
func NewDocker(logger *dockerutil.ContainerLogConsumer) (*dockerutil.DockerManager, error) {
	dockerMu.Lock()
	defer dockerMu.Unlock()

	var initErr error
	dockerOnce.Do(func() {
		cfg := dockerutil.Config{
			ProjectName:      projectName,
			ComposeFile:     "docker-compose.yml",
			RemoteRepo:      remote,
			RequiredServices: []string{"front", "back", "fonts"},
			LogConsumer:     logger,
		}

		dockerInstance, initErr = dockerutil.NewDockerManager(cfg)
	})

	return dockerInstance, initErr
}

// NewAksharamukha creates a new instance of the Aksharamukha service
func NewAksharamukha() (*Aksharamukha, error) {
	logConfig := dockerutil.LogConfig{
		Prefix:      projectName,
		ShowService: true,
		ShowType:    true,
		LogLevel:    zerolog.Disabled,
		InitMessage: "Listening at: http://0.0.0.0:8085",
	}
	
	logger := dockerutil.NewContainerLogConsumer(logConfig)

	docker, err := NewDocker(logger)
	if err != nil {
		return nil, err
	}

	return &Aksharamukha{
		docker: docker,
		logger: logger,
	}, nil
}

// Init initializes the aksharamukha service
func (a *Aksharamukha) Init() error {
	return a.docker.Init()
}

// InitQuiet initializes the aksharamukha service with reduced logging
func (a *Aksharamukha) InitQuiet() error {
	return a.docker.InitQuiet()
}

// InitForce initializes the aksharamukha service with forced rebuild
func (a *Aksharamukha) InitForce() error {
	return a.docker.InitForce()
}

// Stop stops the aksharamukha service
func (a *Aksharamukha) Stop() error {
	return a.docker.Stop()
}

// Close implements io.Closer
func (a *Aksharamukha) Close() error {
	a.logger.Close()
	return a.docker.Close()
}

// Status returns the current status of the aksharamukha service
func (a *Aksharamukha) Status() (string, error) {
	return a.docker.Status()
}

// SetLogLevel updates the logging level
func (a *Aksharamukha) SetLogLevel(level zerolog.Level) {
	a.logger.SetLogLevel(level)
}

func placeholder3456543() {
	color.Redln(" 𝒻*** 𝓎ℴ𝓊 𝒸ℴ𝓂𝓅𝒾𝓁ℯ𝓇")
	pp.Println("𝓯*** 𝔂𝓸𝓾 𝓬𝓸𝓶𝓹𝓲𝓵𝓮𝓻")
}
