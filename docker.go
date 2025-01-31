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
	instance *docker
	once sync.Once
	mu sync.Mutex
)

type docker struct {
	docker *dockerutil.DockerManager
	logger *dockerutil.ContainerLogConsumer
}

// NewDocker creates or returns an existing docker instance
func newDocker() (*docker, error) {
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

		instance = &docker{
			docker: manager,
			logger: logger,
		}
	})

	if initErr != nil {
		return nil, initErr
	}
	return instance, nil
}

// Init initializes the docker service
func Init() error {
	if instance == nil {
		if _, err := newDocker(); err != nil {
			return err
		}
	}
	return instance.docker.Init()
}


// InitQuiet initializes the docker service with reduced logging
func InitQuiet() error {
	if instance == nil {
		if _, err := newDocker(); err != nil {
			return err
		}
	}
	return instance.docker.InitQuiet()
}
// InitRecreate remove existing containers (if noCache is true, downloads the lastest
// version of dependencies ignoring cache), then builds and up the containers
func InitRecreate(noCache bool) error {
	if instance == nil {
		if _, err := newDocker(); err != nil {
			return err
		}
	}
	if noCache {
		return instance.docker.InitRecreateNoCache()
	}
	return instance.docker.InitRecreate()
}

func MustInit() {
	if instance == nil {
		newDocker()
	}
	instance.docker.InitRecreate()
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
