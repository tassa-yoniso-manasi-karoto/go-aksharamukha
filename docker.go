package aksharamukha

import (
	"time"
	"sync"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/gookit/color"
	"github.com/k0kubun/pp"
	
	"github.com/tassa-yoniso-manasi-karoto/dockerutil"
)

const (
	remote = "https://github.com/virtualvinodh/aksharamukha.git"
	projectName = "aksharamukha"
	containerFront = "aksharamukha-front-1"
	containerBack = "aksharamukha-back-1"
	containerFonts = "aksharamukha-fonts-1"
)

var (
	QueryTO = 1 * time.Hour
	instance *Docker
	once sync.Once
	mu sync.Mutex
)

type Docker struct {
	docker *dockerutil.DockerManager
	logger *dockerutil.ContainerLogConsumer
}

// NewDocker creates or returns an existing Docker instance
func NewDocker() (*Docker, error) {
	mu.Lock()
	defer mu.Unlock()
	
	var initErr error
	once.Do(func() {
		logConfig := dockerutil.LogConfig{
			Prefix:      projectName,
			ShowService: true,
			ShowType:    true,
			LogLevel:    zerolog.Disabled, //InfoLevel
			InitMessage: "Listening at: http://0.0.0.0:8085",
		}
		
		logger := dockerutil.NewContainerLogConsumer(logConfig)

		cfg := dockerutil.Config{
			ProjectName:      projectName,
			ComposeFile:     "docker-compose.yml",
			RemoteRepo:      remote,
			RequiredServices: []string{"front", "back", "fonts"},
			LogConsumer:     logger,
		}

		manager, err := dockerutil.NewDockerManager(cfg)
		if err != nil {
			initErr = err
			return
		}

		instance = &Docker{
			docker: manager,
			logger: logger,
		}
	})

	if initErr != nil {
		return nil, initErr
	}
	return instance, nil
}

// Package-level functions for Docker management
func Init() error {
	if instance == nil {
		if _, err := NewDocker(); err != nil {
			return err
		}
	}
	return instance.docker.Init()
}

func InitQuiet() error {
	if instance == nil {
		if _, err := NewDocker(); err != nil {
			return err
		}
	}
	return instance.docker.InitQuiet()
}

func InitForce() error {
	if instance == nil {
		if _, err := NewDocker(); err != nil {
			return err
		}
	}
	return instance.docker.InitForce()
}

func MustInit() {
	if instance == nil {
		NewDocker()
	}
	instance.docker.InitForce()
}


func Stop() error {
	if instance == nil {
		return fmt.Errorf("docker instance not initialized")
	}
	return instance.docker.Stop()
}

func Close() error {
	if instance != nil {
		instance.logger.Close()
		return instance.docker.Close()
	}
	return nil
}

func SetLogLevel(level zerolog.Level) {
	if instance != nil {
		instance.logger.SetLogLevel(level)
	}
}

func placeholder3456543() {
	color.Redln(" 𝒻*** 𝓎ℴ𝓊 𝒸ℴ𝓂𝓅𝒾𝓁ℯ𝓇")
	pp.Println("𝓯*** 𝔂𝓸𝓾 𝓬𝓸𝓶𝓹𝓲𝓵𝓮𝓻")
}
