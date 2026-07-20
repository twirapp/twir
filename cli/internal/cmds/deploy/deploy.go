package deploy

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/urfave/cli/v2"
)

const (
	defaultStackName = "twir"
	defaultRegistry  = "registry.twir.app"
)

var serviceImages = map[string]string{
	"migrations":     "migrations",
	"api-gql":        "api-gql",
	"bots":           "bots",
	"parser":         "parser",
	"timers":         "timers",
	"scheduler":      "scheduler",
	"eventsub":       "eventsub",
	"integrations":   "integrations",
	"web":            "web",
	"overlays":       "overlays",
	"websockets":     "websockets",
	"tokens":         "tokens",
	"emotes-cacher":  "emotes-cacher",
	"events":         "events",
	"dota":           "dota",
	"deploy-webhook": "deploy-receiver",
	"executron":      "executron",
}

var releaseServices = []string{
	"migrations",
	"api-gql",
	"bots",
	"parser",
	"timers",
	"scheduler",
	"eventsub",
	"integrations",
	"web",
	"overlays",
	"websockets",
	"tokens",
	"emotes-cacher",
	"events",
	"dota",
	"executron",
}

type deployConfig struct {
	StackName         string
	Registry          string
	ReceiverAddr      string
	WebhookToken      string
	AllowedTagPattern string
	RequestTimeout    time.Duration
	ApplyTimeout      time.Duration
}

type deployRequest struct {
	Service  string   `json:"service,omitempty"`
	Services []string `json:"services,omitempty"`
	ImageTag string   `json:"imageTag"`
	RefName  string   `json:"refName,omitempty"`
	Trigger  string   `json:"trigger,omitempty"`
}

type deployResult struct {
	Service string `json:"service"`
	Image   string `json:"image"`
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
}

type deployStatus struct {
	Running          bool           `json:"running"`
	CurrentImageTag  string         `json:"currentImageTag,omitempty"`
	LastRequestedTag string         `json:"lastRequestedTag,omitempty"`
	LastRequestedAt  time.Time      `json:"lastRequestedAt"`
	LastStartedAt    time.Time      `json:"lastStartedAt"`
	LastFinishedAt   time.Time      `json:"lastFinishedAt"`
	LastSucceededAt  time.Time      `json:"lastSucceededAt"`
	LastError        string         `json:"lastError,omitempty"`
	LastRefName      string         `json:"lastRefName,omitempty"`
	LastTrigger      string         `json:"lastTrigger,omitempty"`
	Results          []deployResult `json:"results,omitempty"`
}

type receiverServer struct {
	cfg      deployConfig
	tagRegex *regexp.Regexp

	mu     sync.Mutex
	status deployStatus
}

var WebhookCmd = &cli.Command{
	Name:  "deploy-webhook",
	Usage: "start webhook receiver for Docker Swarm service image updates",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "listen-addr", EnvVars: []string{"TWIR_DEPLOY_RECEIVER_ADDR"}, Value: ":8090"},
		&cli.StringFlag{Name: "webhook-token", EnvVars: []string{"TWIR_DEPLOY_WEBHOOK_TOKEN"}, Required: true},
		&cli.StringFlag{Name: "stack-name", EnvVars: []string{"TWIR_DEPLOY_STACK_NAME"}, Value: defaultStackName},
		&cli.StringFlag{Name: "registry", EnvVars: []string{"TWIR_DEPLOY_REGISTRY"}, Value: defaultRegistry},
		&cli.StringFlag{Name: "allowed-tag-pattern", EnvVars: []string{"TWIR_DEPLOY_ALLOWED_TAG_PATTERN"}, Value: `^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`},
		&cli.DurationFlag{Name: "request-timeout", EnvVars: []string{"TWIR_DEPLOY_REQUEST_TIMEOUT"}, Value: 15 * time.Second},
		&cli.DurationFlag{Name: "apply-timeout", EnvVars: []string{"TWIR_DEPLOY_APPLY_TIMEOUT"}, Value: 10 * time.Minute},
	},
	Action: func(c *cli.Context) error {
		cfg := deployConfig{
			StackName:         c.String("stack-name"),
			Registry:          c.String("registry"),
			ReceiverAddr:      c.String("listen-addr"),
			WebhookToken:      strings.TrimSpace(c.String("webhook-token")),
			AllowedTagPattern: c.String("allowed-tag-pattern"),
			RequestTimeout:    c.Duration("request-timeout"),
			ApplyTimeout:      c.Duration("apply-timeout"),
		}

		tagRegex, err := regexp.Compile(cfg.AllowedTagPattern)
		if err != nil {
			return fmt.Errorf("compile allowed tag pattern: %w", err)
		}

		server := &receiverServer{cfg: cfg, tagRegex: tagRegex}

		mux := http.NewServeMux()
		mux.HandleFunc("/healthz", server.handleHealthz)
		mux.HandleFunc("/status", server.handleStatus)
		mux.HandleFunc("/deploy", server.handleDeploy)

		httpServer := &http.Server{
			Addr:              cfg.ReceiverAddr,
			Handler:           mux,
			ReadHeaderTimeout: cfg.RequestTimeout,
			ReadTimeout:       cfg.RequestTimeout,
			WriteTimeout:      cfg.RequestTimeout,
			IdleTimeout:       60 * time.Second,
		}

		log.Printf("deploy webhook listening on %s", cfg.ReceiverAddr)

		return httpServer.ListenAndServe()
	},
}

var ApplyCmd = &cli.Command{
	Name:   "deploy-apply",
	Usage:  "update one Swarm service image tag",
	Hidden: true,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "service", Required: true},
		&cli.StringFlag{Name: "image-tag", Required: true},
		&cli.StringFlag{Name: "stack-name", EnvVars: []string{"TWIR_DEPLOY_STACK_NAME"}, Value: defaultStackName},
		&cli.StringFlag{Name: "registry", EnvVars: []string{"TWIR_DEPLOY_REGISTRY"}, Value: defaultRegistry},
		&cli.DurationFlag{Name: "apply-timeout", EnvVars: []string{"TWIR_DEPLOY_APPLY_TIMEOUT"}, Value: 10 * time.Minute},
	},
	Action: func(c *cli.Context) error {
		cfg := deployConfig{
			StackName:    c.String("stack-name"),
			Registry:     c.String("registry"),
			ApplyTimeout: c.Duration("apply-timeout"),
		}

		result, err := updateServiceImage(c.Context, cfg, c.String("service"), c.String("image-tag"))
		if err != nil {
			return err
		}

		payload, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(payload))

		return nil
	},
}

func (s *receiverServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *receiverServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	status := s.status
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, status)
}

func (s *receiverServer) handleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.authorized(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)
	defer r.Body.Close()

	var req deployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("decode request: %v", err), http.StatusBadRequest)
		return
	}

	req.ImageTag = strings.TrimSpace(req.ImageTag)
	if req.ImageTag == "" {
		http.Error(w, "imageTag is required", http.StatusBadRequest)
		return
	}
	if !s.tagRegex.MatchString(req.ImageTag) {
		http.Error(w, "imageTag does not match allowed pattern", http.StatusBadRequest)
		return
	}

	services := servicesFromRequest(req)
	if len(services) == 0 {
		http.Error(w, "service or services is required", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	if s.status.Running {
		s.mu.Unlock()
		http.Error(w, "deploy already running", http.StatusConflict)
		return
	}

	now := time.Now()
	s.status.Running = true
	s.status.LastRequestedAt = now
	s.status.LastRequestedTag = req.ImageTag
	s.status.LastRefName = req.RefName
	s.status.LastTrigger = req.Trigger
	s.status.LastError = ""
	s.status.Results = nil
	s.mu.Unlock()

	go s.runDeploy(req, services)

	writeJSON(w, http.StatusAccepted, map[string]any{
		"status":   "accepted",
		"imageTag": req.ImageTag,
		"services": services,
	})
}

func (s *receiverServer) runDeploy(req deployRequest, services []string) {
	startedAt := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ApplyTimeout)
	defer cancel()

	results := make([]deployResult, 0, len(services))
	var deployErr error

	for _, service := range services {
		result, err := updateServiceImage(ctx, s.cfg, service, req.ImageTag)
		if err != nil {
			result = deployResult{Service: service, Status: "failed", Error: err.Error()}
			deployErr = errors.Join(deployErr, fmt.Errorf("%s: %w", service, err))
		}

		results = append(results, result)
	}

	s.finishDeploy(req.ImageTag, startedAt, results, deployErr)
}

func (s *receiverServer) finishDeploy(imageTag string, startedAt time.Time, results []deployResult, err error) {
	finishedAt := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.Running = false
	s.status.LastStartedAt = startedAt
	s.status.LastFinishedAt = finishedAt
	s.status.Results = results

	if err != nil {
		s.status.LastError = err.Error()
		log.Printf("deploy failed for tag %s: %v", imageTag, err)
		return
	}

	s.status.CurrentImageTag = imageTag
	s.status.LastSucceededAt = finishedAt
	s.status.LastError = ""
	log.Printf("deploy finished for tag %s", imageTag)
}

func (s *receiverServer) authorized(r *http.Request) bool {
	provided := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	if provided == "" {
		provided = strings.TrimSpace(r.Header.Get("X-Deploy-Token"))
	}

	return subtle.ConstantTimeCompare([]byte(provided), []byte(s.cfg.WebhookToken)) == 1
}

func servicesFromRequest(req deployRequest) []string {
	services := make([]string, 0, len(req.Services)+1)

	if strings.TrimSpace(req.Service) != "" {
		services = append(services, strings.TrimSpace(req.Service))
	}

	for _, service := range req.Services {
		service = strings.TrimSpace(service)
		if service == "" || slices.Contains(services, service) {
			continue
		}

		services = append(services, service)
	}

	if len(services) == 0 && req.Trigger == "github-actions" {
		return slices.Clone(releaseServices)
	}

	return services
}

func updateServiceImage(ctx context.Context, cfg deployConfig, serviceName string, imageTag string) (deployResult, error) {
	imageRepo, ok := serviceImages[serviceName]
	if !ok {
		return deployResult{Service: serviceName, Status: "failed"}, fmt.Errorf("unknown managed service %q", serviceName)
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return deployResult{Service: serviceName, Status: "failed"}, fmt.Errorf("create docker client: %w", err)
	}
	defer dockerClient.Close()

	service, err := findService(ctx, dockerClient, cfg.StackName, serviceName)
	if err != nil {
		return deployResult{Service: serviceName, Status: "failed"}, err
	}

	if service.Spec.TaskTemplate.ContainerSpec == nil {
		return deployResult{Service: serviceName, Status: "failed"}, errors.New("service has no container spec")
	}

	image := fmt.Sprintf("%s/twirapp/%s:%s", cfg.Registry, imageRepo, imageTag)
	if service.Spec.TaskTemplate.ContainerSpec.Image == image {
		return deployResult{Service: serviceName, Image: image, Status: "unchanged"}, nil
	}

	spec := service.Spec
	spec.TaskTemplate.ContainerSpec.Image = image
	if spec.Labels == nil {
		spec.Labels = make(map[string]string)
	}
	spec.Labels["twir.deploy.image-tag"] = imageTag
	spec.Labels["twir.deploy.updated-at"] = time.Now().UTC().Format(time.RFC3339)

	_, err = dockerClient.ServiceUpdate(ctx, service.ID, service.Version, spec, types.ServiceUpdateOptions{})
	if err != nil {
		return deployResult{Service: serviceName, Image: image, Status: "failed"}, fmt.Errorf("update swarm service: %w", err)
	}

	return deployResult{Service: serviceName, Image: image, Status: "updated"}, nil
}

func findService(ctx context.Context, dockerClient *client.Client, stackName string, serviceName string) (swarm.Service, error) {
	services, err := dockerClient.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		return swarm.Service{}, fmt.Errorf("list swarm services: %w", err)
	}

	fullServiceName := stackName + "_" + serviceName
	for _, service := range services {
		if service.Spec.Name == fullServiceName && service.Spec.Labels["com.docker.stack.namespace"] == stackName {
			return service, nil
		}
	}

	return swarm.Service{}, fmt.Errorf("service %q not found in stack %q", serviceName, stackName)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
