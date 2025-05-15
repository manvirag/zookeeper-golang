package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger = logrus.New()

type Service struct {
	consul    *consulapi.Client
	port      int
	server    *http.Server
	router    *mux.Router
	serviceID string
}

func NewService(consulHost string, serverId string, port int) (*Service, error) {
	// Create Consul client config
	config := consulapi.DefaultConfig()
	config.Address = consulHost

	// Create Consul client
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %v", err)
	}

	router := mux.NewRouter()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return &Service{
		consul:    client,
		port:      port,
		server:    server,
		router:    router,
		serviceID: serverId,
	}, nil
}

func (s *Service) registerService(serviceName string) error {
	registration := &consulapi.AgentServiceRegistration{
		ID:      s.serviceID,
		Name:    serviceName,
		Address: s.serviceID,
		Port:    s.port,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", s.serviceID, s.port),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	if err := s.consul.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	logger.Infof("Service registered with ID: %s", s.serviceID)
	return nil
}

func (s *Service) setupRoutes() {
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Service running on port %d", s.port)))
	}).Methods("GET")
}

func (s *Service) Start() error {
	s.setupRoutes()

	// Get service name from environment variable
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "my-service"
	}

	// Register service with Consul
	if err := s.registerService(serviceName); err != nil {
		return err
	}

	// Start HTTP server
	go func() {
		logger.Infof("Starting service on port %d", s.port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	return nil
}

func (s *Service) Stop() error {
	// Deregister service from Consul
	if err := s.consul.Agent().ServiceDeregister(s.serviceID); err != nil {
		logger.Errorf("Failed to deregister service: %v", err)
	}

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func main() {
	// Get configuration from environment variables
	port := 8080
	if portStr := os.Getenv("SERVICE_PORT"); portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
	}

	consulHost := os.Getenv("CONSUL_HOST")
	if consulHost == "" {
		consulHost = "consul:8500"
	}

	serverId := os.Getenv("SERVER_ID")

	if serverId == "" {
		return
	}

	service, err := NewService(consulHost, serverId, port)
	if err != nil {
		logger.Fatalf("Failed to create service: %v", err)
	}

	if err := service.Start(); err != nil {
		logger.Fatalf("Failed to start service: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down service...")
	if err := service.Stop(); err != nil {
		logger.Fatalf("Failed to stop service: %v", err)
	}
}
