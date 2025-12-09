package aksharamukha

import (
	"time"
	"sync"
	"context"
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

	// Docker Hub images
	imageFront = "virtualvinodh/aksharamukha-front"
	imageBack  = "virtualvinodh/aksharamukha-back"
	imageFonts = "virtualvinodh/aksharamukha-fonts"
)

var (
	DefaultQueryTimeout = 5 * time.Minute
	DefaultDockerLogLevel = zerolog.TraceLevel
)

// AksharamukhaManager handles Docker lifecycle for Aksharamukha project
type AksharamukhaManager struct {
	docker       *dockerutil.DockerManager
	logger       *dockerutil.ContainerLogConsumer
	projectName  string
	frontContainer string
	backContainer  string
	fontsContainer string
	QueryTimeout time.Duration
}

// ManagerOption defines function signature for options to configure AksharamukhaManager
type ManagerOption func(*AksharamukhaManager)

// WithQueryTimeout sets a custom query timeout
func WithQueryTimeout(timeout time.Duration) ManagerOption {
	return func(am *AksharamukhaManager) {
		am.QueryTimeout = timeout
	}
}

// WithProjectName sets a custom project name for multiple instances
func WithProjectName(name string) ManagerOption {
	return func(am *AksharamukhaManager) {
		am.projectName = name
		// Default container names are derived from project name
		am.frontContainer = name + "-front-1"
		am.backContainer = name + "-back-1"
		am.fontsContainer = name + "-fonts-1"
	}
}

// WithContainerNames overrides the default container names
func WithContainerNames(front, back, fonts string) ManagerOption {
	return func(am *AksharamukhaManager) {
		am.frontContainer = front
		am.backContainer = back
		am.fontsContainer = fonts
	}
}

// NewManager creates a new Aksharamukha manager instance
func NewManager(ctx context.Context, opts ...ManagerOption) (*AksharamukhaManager, error) {
	manager := &AksharamukhaManager{
		projectName: projectName,
		frontContainer: containerFront,
		backContainer: containerBack,
		fontsContainer: containerFonts,
		QueryTimeout: DefaultQueryTimeout,
	}
	
	// Apply options
	for _, opt := range opts {
		opt(manager)
	}
	
	logConfig := dockerutil.LogConfig{
		Prefix:      manager.projectName,
		ShowService: true,
		ShowType:    true,
		LogLevel:    DefaultDockerLogLevel,
		InitMessage: "Listening at: http://0.0.0.0:8085",
	}

	logger := dockerutil.NewContainerLogConsumer(logConfig)

	cfg := dockerutil.Config{
		ProjectName:      manager.projectName,
		ComposeFile:      "docker-compose.yml",
		RemoteRepo:       remote,
		RequiredServices: []string{"front", "back", "fonts"},
		LogConsumer:      logger,
		Timeout: dockerutil.Timeout{
			Create:   60 * time.Second,
			Recreate: 10 * time.Minute,
			Start:    60 * time.Second,
		},
	}

	dockerManager, err := dockerutil.NewDockerManager(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker manager: %w", err)
	}

	manager.docker = dockerManager
	manager.logger = logger
	
	return manager, nil
}

// Init initializes the docker service
func (am *AksharamukhaManager) Init(ctx context.Context) error {
	return am.docker.Init()
}

// InitQuiet initializes the docker service with reduced logging
func (am *AksharamukhaManager) InitQuiet(ctx context.Context) error {
	return am.docker.InitQuiet()
}

// InitRecreate remove existing containers then builds and up the containers
func (am *AksharamukhaManager) InitRecreate(ctx context.Context, noCache bool) error {
	if noCache {
		return am.docker.InitRecreateNoCache()
	}
	return am.docker.InitRecreate()
}

// PullImages pre-pulls all required Docker images with retry logic.
// This is useful for slow/unreliable connections as it provides better
// error handling than docker-compose's built-in pull.
func (am *AksharamukhaManager) PullImages(ctx context.Context) error {
	images := []string{imageFront, imageBack, imageFonts}
	opts := dockerutil.DefaultPullOptions()

	for _, img := range images {
		if err := dockerutil.PullImage(ctx, img, opts); err != nil {
			return fmt.Errorf("failed to pull image %s: %w", img, err)
		}
	}
	return nil
}

// MustInit initializes the docker service and panics on error
func (am *AksharamukhaManager) MustInit(ctx context.Context) {
	if err := am.docker.InitRecreate(); err != nil {
		panic(err)
	}
}

// Stop stops the docker service
func (am *AksharamukhaManager) Stop(ctx context.Context) error {
	return am.docker.Stop()
}

// Close implements io.Closer
func (am *AksharamukhaManager) Close() error {
	am.logger.Close()
	return am.docker.Close()
}

// GetBaseURL returns the base URL for API requests
func (am *AksharamukhaManager) GetBaseURL() string {
	return "http://localhost:8085/api/public"
}

// For backward compatibility with existing code
var (
	instance *AksharamukhaManager
	mu sync.Mutex
	instanceClosed bool
)

// InitWithContext initializes the default docker service with a context
func InitWithContext(ctx context.Context) error {
	mgr, err := getOrCreateDefaultManager(ctx)
	if err != nil {
		return err
	}
	return mgr.Init(ctx)
}

// Init initializes the default docker service (backward compatibility)
func Init() error {
	return InitWithContext(context.Background())
}

// InitQuietWithContext initializes the docker service with reduced logging and a context
func InitQuietWithContext(ctx context.Context) error {
	mgr, err := getOrCreateDefaultManager(ctx)
	if err != nil {
		return err
	}
	return mgr.InitQuiet(ctx)
}

// InitQuiet initializes the docker service with reduced logging (backward compatibility)
func InitQuiet() error {
	return InitQuietWithContext(context.Background())
}

// InitRecreateWithContext removes existing containers with a context
func InitRecreateWithContext(ctx context.Context, noCache bool) error {
	mgr, err := getOrCreateDefaultManager(ctx)
	if err != nil {
		return err
	}
	return mgr.InitRecreate(ctx, noCache)
}

// InitRecreate removes existing containers (backward compatibility)
func InitRecreate(noCache bool) error {
	return InitRecreateWithContext(context.Background(), noCache)
}

// PullImagesWithContext pre-pulls all required Docker images with retry logic
func PullImagesWithContext(ctx context.Context) error {
	mgr, err := getOrCreateDefaultManager(ctx)
	if err != nil {
		return err
	}
	return mgr.PullImages(ctx)
}

// PullImages pre-pulls all required Docker images (backward compatibility)
func PullImages() error {
	return PullImagesWithContext(context.Background())
}

// MustInitWithContext initializes the docker service with a context (panics on error)
func MustInitWithContext(ctx context.Context) {
	mgr, _ := getOrCreateDefaultManager(ctx)
	mgr.MustInit(ctx)
}

// MustInit initializes the docker service (backward compatibility)
func MustInit() {
	MustInitWithContext(context.Background())
}

// StopWithContext stops the docker service with a context
func StopWithContext(ctx context.Context) error {
	if instance == nil {
		return fmt.Errorf("docker instance not initialized")
	}
	return instance.Stop(ctx)
}

// Stop stops the docker service (backward compatibility)
func Stop() error {
	return StopWithContext(context.Background())
}

// Close implements io.Closer (backward compatibility)
func Close() error {
	mu.Lock()
	defer mu.Unlock()
	
	if instance != nil {
		instance.logger.Close()
		err := instance.docker.Close()
		// Mark the instance as closed
		instanceClosed = true
		return err
	}
	return nil
}

// getOrCreateDefaultManager returns or creates the default manager instance
func getOrCreateDefaultManager(ctx context.Context) (*AksharamukhaManager, error) {
	mu.Lock()
	defer mu.Unlock()
	
	// Create a new instance if it doesn't exist or was previously closed
	if instance == nil || instanceClosed {
		mgr, err := NewManager(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create default manager: %w", err)
		}
		instance = mgr
		instanceClosed = false
	}
	
	return instance, nil
}

func placeholder3456543() {
	color.Redln(" 𝒻*** 𝓎ℴ𝓊 𝒸ℴ𝓂𝓅𝒾𝓁ℯ𝓇")
	pp.Println("𝓯*** 𝔂𝓸𝓾 𝓬𝓸𝓶𝓹𝓲𝓵𝓮𝓻")
}