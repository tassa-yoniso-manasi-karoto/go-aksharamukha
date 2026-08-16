package aksharamukha

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/gookit/color"
	"github.com/k0kubun/pp"
	"github.com/rs/zerolog"

	"github.com/tassa-yoniso-manasi-karoto/dockerutil"
)

const (
	projectName   = "aksharamukha"
	containerBack = "aksharamukha-back-1"
	internalPort  = 8085

	// Docker Hub image for API backend (front/fonts not needed - web UI only)
	imageBack = "virtualvinodh/aksharamukha-back"
)

var (
	DefaultQueryTimeout   = 5 * time.Minute
	DefaultDockerLogLevel = zerolog.TraceLevel

	// Package-level download progress callback for image pulls
	downloadProgressCallback func(current, total int64, status string)
	downloadCallbackMu       sync.Mutex
)

// SetDownloadProgressCallback sets the callback for image download progress.
// This should be called before PullImages/PullImagesWithContext.
func SetDownloadProgressCallback(cb func(current, total int64, status string)) {
	downloadCallbackMu.Lock()
	defer downloadCallbackMu.Unlock()
	downloadProgressCallback = cb
}

// ClearDownloadProgressCallback clears the download progress callback.
func ClearDownloadProgressCallback() {
	downloadCallbackMu.Lock()
	defer downloadCallbackMu.Unlock()
	downloadProgressCallback = nil
}

// AksharamukhaManager handles Docker lifecycle for Aksharamukha project
type AksharamukhaManager struct {
	docker                   *dockerutil.DockerManager
	logger                   *dockerutil.ContainerLogConsumer
	projectName              string
	backContainer            string
	containerNameExplicit    bool
	baseURL                  string
	QueryTimeout             time.Duration
	downloadProgressCallback func(current, total int64, status string)
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
		am.backContainer = name + "-back-1"
	}
}

// WithContainerName overrides the default container name
func WithContainerName(name string) ManagerOption {
	return func(am *AksharamukhaManager) {
		am.backContainer = name
		am.containerNameExplicit = true
	}
}

// WithDownloadProgressCallback sets a callback for download progress during image pull
func WithDownloadProgressCallback(cb func(current, total int64, status string)) ManagerOption {
	return func(am *AksharamukhaManager) {
		am.downloadProgressCallback = cb
	}
}

// buildComposeProject creates the compose project definition for aksharamukha
// Only the "back" service is needed - front/fonts are for the web UI
func buildComposeProject(projectName, containerName string) *types.Project {
	// Network name follows Docker Compose convention: {project}_{network}
	defaultNetworkName := projectName + "_default"

	return &types.Project{
		Name: projectName,
		// Default network required for port exposure
		Networks: types.Networks{
			"default": types.NetworkConfig{
				Name: defaultNetworkName,
			},
		},
		Services: types.Services{
			"back": {
				Name:          "back",
				ContainerName: containerName,
				Image:         imageBack,
				Ports: []types.ServicePortConfig{{
					Published: "0",
					Target:    internalPort,
					Protocol:  "tcp",
					Mode:      "ingress",
				}},
				// Attach to default network
				Networks: map[string]*types.ServiceNetworkConfig{
					"default": nil,
				},
			},
		},
	}
}

// NewManager creates a new Aksharamukha manager instance
func NewManager(ctx context.Context, opts ...ManagerOption) (*AksharamukhaManager, error) {
	manager := &AksharamukhaManager{
		projectName:   projectName,
		backContainer: containerBack,
		QueryTimeout:  DefaultQueryTimeout,
	}

	// Apply options
	for _, opt := range opts {
		opt(manager)
	}

	var projectLease *dockerutil.ProjectLease
	managerReady := false
	defer func() {
		if projectLease != nil && !managerReady {
			_ = projectLease.Release(ctx)
		}
	}()
	if runtime := dockerutil.RuntimeFromContext(ctx); runtime != nil {
		lease, err := runtime.AcquireProject(ctx, dockerutil.ProjectSpec{
			BaseName:  manager.projectName,
			Kind:      "aksharamukha",
			Lifecycle: dockerutil.LifecycleOwned,
		})
		if err != nil {
			return nil, fmt.Errorf("acquire Aksharamukha project: %w", err)
		}
		projectLease = lease
		manager.projectName = lease.Name()
		if !manager.containerNameExplicit {
			manager.backContainer = ""
		}
	}

	// Build compose project
	project := buildComposeProject(manager.projectName, manager.backContainer)

	logConfig := dockerutil.LogConfig{
		Prefix:      manager.projectName,
		ShowService: true,
		ShowType:    true,
		LogLevel:    DefaultDockerLogLevel,
		InitMessage: "Listening at: http://0.0.0.0:8085",
	}

	logger := dockerutil.NewContainerLogConsumer(logConfig)

	cfg := dockerutil.Config{
		Project:          project,
		ProjectLease:     projectLease,
		RequiredServices: []string{"back"},
		LogConsumer:      logger,
		Timeout: dockerutil.Timeout{
			Create:   60 * time.Second,
			Recreate: 10 * time.Minute,
			Start:    60 * time.Second,
		},
		OnPullProgress: manager.downloadProgressCallback,
	}

	dockerManager, err := dockerutil.NewDockerManager(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker manager: %w", err)
	}

	manager.docker = dockerManager
	manager.logger = logger
	managerReady = true

	return manager, nil
}

// Init initializes the docker service
func (am *AksharamukhaManager) Init(ctx context.Context) error {
	if err := am.docker.Init(); err != nil {
		return err
	}
	return am.resolveServiceEndpoint(ctx)
}

// InitQuiet initializes the docker service with reduced logging
func (am *AksharamukhaManager) InitQuiet(ctx context.Context) error {
	if err := am.docker.InitQuiet(); err != nil {
		return err
	}
	return am.resolveServiceEndpoint(ctx)
}

// InitRecreate remove existing containers then builds and up the containers
func (am *AksharamukhaManager) InitRecreate(ctx context.Context, noCache bool) error {
	if noCache {
		if err := am.docker.InitRecreateNoCache(); err != nil {
			return err
		}
	} else if err := am.docker.InitRecreate(); err != nil {
		return err
	}
	return am.resolveServiceEndpoint(ctx)
}

func (am *AksharamukhaManager) resolveServiceEndpoint(ctx context.Context) error {
	containerID, err := am.docker.ContainerID(ctx, "back")
	if err != nil {
		return fmt.Errorf("resolve Aksharamukha container: %w", err)
	}
	host, port, err := am.docker.PublishedPort(ctx, "back", internalPort)
	if err != nil {
		return fmt.Errorf("resolve Aksharamukha service port: %w", err)
	}
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "localhost"
	}
	am.backContainer = containerID
	am.baseURL = fmt.Sprintf("http://%s:%d/api/public", host, port)
	return nil
}

// PullImages pre-pulls all required Docker images with retry logic.
// This is useful for slow/unreliable connections as it provides better
// error handling than docker-compose's built-in pull.
func (am *AksharamukhaManager) PullImages(ctx context.Context) error {
	images := []string{imageBack}

	opts := dockerutil.DefaultPullOptions()

	// Use manager's callback if set, otherwise fall back to package-level
	if am.downloadProgressCallback != nil {
		opts.OnProgress = am.downloadProgressCallback
	} else {
		downloadCallbackMu.Lock()
		cb := downloadProgressCallback
		downloadCallbackMu.Unlock()
		if cb != nil {
			opts.OnProgress = cb
		}
	}

	// dockerutil.PullImages handles:
	// - Manifest fetching for accurate total size
	// - Layer deduplication across images
	// - Unified progress tracking
	return dockerutil.PullImages(ctx, images, opts)
}

// MustInit initializes the docker service and panics on error
func (am *AksharamukhaManager) MustInit(ctx context.Context) {
	if err := am.InitRecreate(ctx, false); err != nil {
		panic(err)
	}
}

// Stop stops the docker service
func (am *AksharamukhaManager) Stop(ctx context.Context) error {
	return am.docker.StopWithContext(ctx)
}

// Close implements io.Closer
func (am *AksharamukhaManager) Close() error {
	return am.CloseWithContext(context.Background())
}

// CloseWithContext releases the manager's Docker ownership with a context.
func (am *AksharamukhaManager) CloseWithContext(ctx context.Context) error {
	am.logger.Close()
	return am.docker.CloseWithContext(ctx)
}

// GetBaseURL returns the base URL for API requests
func (am *AksharamukhaManager) GetBaseURL() string {
	return am.baseURL
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
		err := instance.docker.CloseWithContext(context.Background())
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
