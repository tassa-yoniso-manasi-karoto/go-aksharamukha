package aksharamukha

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"sync"
	"time"

	"github.com/adrg/xdg"
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

	// Logger for this package
	Logger = zerolog.Nop()

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
	runtime                  *dockerutil.Runtime // the manager's own, when it was made without one; closed with it
	logger                   *dockerutil.ContainerLogConsumer
	projectName              string
	backContainer            string
	containerNameExplicit    bool
	endpointMu               sync.RWMutex // guards backContainer, baseURL and resolution, which recovery replaces while calls read them
	baseURL                  string
	resolution               int        // counts the resolutions of baseURL, so that calls that failed on the same one recover once
	recoverMu                sync.Mutex // one call at a time brings the service back
	closed                   bool       // under recoverMu: a closed manager brings nothing back
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
		if manager.runtime != nil && !managerReady {
			_ = manager.runtime.Close(ctx)
		}
	}()
	runtime := dockerutil.RuntimeFromContext(ctx)
	if runtime == nil {
		// Without a runtime, every process shared one unowned project, and the
		// first to close removed it under the others, whose calls then failed on
		// its port. A runtime of the manager's own gives it a project that only
		// it removes, and reaps those that crashed processes left.
		private, err := dockerutil.NewRuntime(ctx, dockerutil.RuntimeConfig{
			Application: "go-aksharamukha",
			RootDir:     filepath.Join(xdg.CacheHome, "go-aksharamukha", "docker"),
		})
		if err != nil {
			return nil, fmt.Errorf("create Aksharamukha Docker runtime: %w", err)
		}
		if err := private.RecoveryError(); err != nil {
			Logger.Warn().Err(err).Msg("Removing the containers of crashed processes was deferred")
		}
		manager.runtime = private
		runtime = private
	}
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
	am.endpointMu.Lock()
	am.backContainer = containerID
	am.baseURL = fmt.Sprintf("http://%s:%d/api/public", host, port)
	am.resolution++
	am.endpointMu.Unlock()
	return nil
}

// endpoint returns the service's base URL and the resolution it comes from.
func (am *AksharamukhaManager) endpoint() (baseURL string, resolution int) {
	am.endpointMu.RLock()
	defer am.endpointMu.RUnlock()
	return am.baseURL, am.resolution
}

// call runs request against the service's base URL. When the service
// cannot be reached, as when its container was stopped or removed, or
// Docker restarted, which also gives it a new port, it brings the
// container back and runs request once more.
func (am *AksharamukhaManager) call(ctx context.Context, request func(baseURL string) error) error {
	baseURL, resolution := am.endpoint()
	err := request(baseURL)
	if err == nil || !unreachable(err) {
		return err
	}
	if recoverErr := am.recover(ctx, resolution); recoverErr != nil {
		return fmt.Errorf("%w; bringing the Aksharamukha service back failed: %v", err, recoverErr)
	}
	baseURL, _ = am.endpoint()
	err = request(baseURL)
	if unreachable(err) {
		// The container listens again, so something between it and this
		// host, such as a VPN, keeps the call from it.
		return fmt.Errorf("the Aksharamukha service is still unreachable after its container restarted (THIS ERROR MAY BE CAUSED BY AN ACTIVE VPN, STOP IT AND RESTART %s): %w",
			dockerutil.DockerBackendName(), err)
	}
	return err
}

// recover restarts the container and finds its endpoint again, after a
// call could not reach the service at the endpoint of resolution failed.
// Calls that failed together recover once: when another one already
// resolved the endpoint anew, the caller only retries.
func (am *AksharamukhaManager) recover(ctx context.Context, failed int) error {
	am.recoverMu.Lock()
	defer am.recoverMu.Unlock()
	if am.closed {
		return errors.New("the manager is closed")
	}
	if _, resolution := am.endpoint(); resolution != failed {
		return nil
	}
	Logger.Warn().Str("project", am.projectName).Msg("Aksharamukha service unreachable; restarting its container")
	if err := am.docker.Restart(ctx, false); err != nil {
		return err
	}
	// The restart goes on even when ctx ends, and so must finding the new
	// port, or the next call would restart the container again.
	return am.resolveServiceEndpoint(context.WithoutCancel(ctx))
}

// unreachable reports whether err means that nothing answered at the
// service's address, or that the connection ended before an answer,
// rather than an answer or a timeout.
func unreachable(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) && opErr.Op == "dial" {
		return true
	}
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)
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
	// A call that fails on the removed container must not bring it back.
	am.recoverMu.Lock()
	am.closed = true
	am.recoverMu.Unlock()

	am.logger.Close()
	err := am.docker.CloseWithContext(ctx)
	if am.runtime != nil {
		err = errors.Join(err, am.runtime.Close(ctx))
	}
	return err
}

// GetBaseURL returns the base URL for API requests
func (am *AksharamukhaManager) GetBaseURL() string {
	baseURL, _ := am.endpoint()
	return baseURL
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
		err := instance.Close()
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
