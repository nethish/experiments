package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
)

type App struct {
	serviceName  string
	servicePort  int
	consulURL    string
	consulClient *api.Client
	serviceID    string
}

type ServiceInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
	Status  string `json:"status"`
}

func NewApp() (*App, error) {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "default-app"
	}

	servicePortStr := os.Getenv("SERVICE_PORT")
	if servicePortStr == "" {
		servicePortStr = "8080"
	}
	servicePort, err := strconv.Atoi(servicePortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVICE_PORT: %v", err)
	}

	consulURL := os.Getenv("CONSUL_URL")
	if consulURL == "" {
		consulURL = "localhost:8500"
	}

	// Create Consul client
	config := api.DefaultConfig()
	config.Address = consulURL
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %v", err)
	}

	serviceID := fmt.Sprintf("%s-%d", serviceName, time.Now().Unix())

	return &App{
		serviceName:  serviceName,
		servicePort:  servicePort,
		consulURL:    consulURL,
		consulClient: client,
		serviceID:    serviceID,
	}, nil
}

func (a *App) registerService() error {
	// Get container hostname (Docker container name)
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %v", err)
	}

	// Register service with Consul
	registration := &api.AgentServiceRegistration{
		ID:      a.serviceID,
		Name:    a.serviceName,
		Port:    a.servicePort,
		Address: hostname,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", hostname, a.servicePort),
			Interval: "10s",
			Timeout:  "3s",
		},
	}

	err = a.consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	log.Printf("Service %s registered with ID %s at %s:%d", a.serviceName, a.serviceID, hostname, a.servicePort)
	return nil
}

func (a *App) deregisterService() error {
	err := a.consulClient.Agent().ServiceDeregister(a.serviceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %v", err)
	}
	log.Printf("Service %s deregistered", a.serviceID)
	return nil
}

func (a *App) discoverServices() ([]ServiceInfo, error) {
	services, _, err := a.consulClient.Health().State("any", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query services: %v", err)
	}

	var serviceInfos []ServiceInfo
	serviceMap := make(map[string]ServiceInfo)

	for _, service := range services {
		if service.ServiceName == "" {
			continue
		}

		// Get service details
		serviceInfo := ServiceInfo{
			ID:      service.ServiceID,
			Name:    service.ServiceName,
			Address: service.Node,
			Status:  service.Status,
		}

		// Get port from service registration
		agentServices, err := a.consulClient.Agent().Services()
		if err == nil {
			if agentService, exists := agentServices[service.ServiceID]; exists {
				serviceInfo.Port = agentService.Port
				serviceInfo.Address = agentService.Address
			}
		}

		serviceMap[service.ServiceID] = serviceInfo
	}

	for _, info := range serviceMap {
		serviceInfos = append(serviceInfos, info)
	}

	return serviceInfos, nil
}

func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (a *App) statusHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service": a.serviceName,
		"id":      a.serviceID,
		"port":    a.servicePort,
		"status":  "running",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (a *App) servicesHandler(w http.ResponseWriter, r *http.Request) {
	services, err := a.discoverServices()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to discover services: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"discovered_services": services,
		"count":               len(services),
	})
}

func (a *App) callServiceHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("service")
	if serviceName == "" {
		http.Error(w, "service parameter is required", http.StatusBadRequest)
		return
	}

	// Discover specific service
	services, _, err := a.consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find service %s: %v", serviceName, err), http.StatusInternalServerError)
		return
	}

	if len(services) == 0 {
		http.Error(w, fmt.Sprintf("Service %s not found", serviceName), http.StatusNotFound)
		return
	}

	// Use first healthy service instance
	service := services[0]
	serviceURL := fmt.Sprintf("http://%s:%d/status", service.Service.Address, service.Service.Port)

	// Make HTTP call to the service
	resp, err := http.Get(serviceURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to call service %s: %v", serviceName, err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"called_service":  serviceName,
		"service_url":     serviceURL,
		"response_status": resp.Status,
		"message":         fmt.Sprintf("Successfully called %s", serviceName),
	})
}

func (a *App) watchServices() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			services, err := a.discoverServices()
			if err != nil {
				log.Printf("Error discovering services: %v", err)
				continue
			}

			log.Printf("Current services in cluster:")
			for _, service := range services {
				log.Printf("  - %s (%s) at %s:%d [Status: %s]",
					service.Name, service.ID, service.Address, service.Port, service.Status)
			}
		}
	}
}

func (a *App) setupRoutes() {
	http.HandleFunc("/health", a.healthHandler)
	http.HandleFunc("/status", a.statusHandler)
	http.HandleFunc("/services", a.servicesHandler)
	http.HandleFunc("/call", a.callServiceHandler)

	// Serve basic info page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := fmt.Sprintf(`
		<html>
		<head><title>%s</title></head>
		<body>
			<h1>Service: %s</h1>
			<p>Service ID: %s</p>
			<p>Port: %d</p>
			<h2>Available Endpoints:</h2>
			<ul>
				<li><a href="/status">/status</a> - Service status</li>
				<li><a href="/services">/services</a> - Discover all services</li>
				<li><a href="/health">/health</a> - Health check</li>
				<li>/call?service=SERVICE_NAME - Call another service</li>
			</ul>
		</body>
		</html>
		`, a.serviceName, a.serviceName, a.serviceID, a.servicePort)

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})
}

func (a *App) run() error {
	// Wait for Consul to be ready
	log.Println("Waiting for Consul to be ready...")
	for i := 0; i < 30; i++ {
		_, err := a.consulClient.Status().Leader()
		if err == nil {
			break
		}
		if i == 29 {
			return fmt.Errorf("consul not ready after 30 attempts")
		}
		time.Sleep(1 * time.Second)
	}

	// Register service
	if err := a.registerService(); err != nil {
		return err
	}

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down...")
		if err := a.deregisterService(); err != nil {
			log.Printf("Error deregistering service: %v", err)
		}
		os.Exit(0)
	}()

	// Start service discovery watcher
	go a.watchServices()

	// Setup HTTP routes
	a.setupRoutes()

	// Start HTTP server
	addr := fmt.Sprintf(":%d", a.servicePort)
	log.Printf("Starting %s on port %d", a.serviceName, a.servicePort)
	log.Printf("Consul UI available at: http://localhost:8500")

	return http.ListenAndServe(addr, nil)
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err := app.run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
