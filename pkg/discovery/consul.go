package discovery

import (
	"fmt"
	"log"

	consul "github.com/hashicorp/consul/api"
)

// Registry defines the interface for service discovery operations (Secondary Port)
type Registry interface {
	Register(serviceID, serviceName, host string, port int) error
	Deregister(serviceID string) error
	Discover(serviceName string) (string, error)
}

// ConsulRegistry is a Consul-backed implementation of the Registry interface (Secondary Adapter)
type ConsulRegistry struct {
	client *consul.Client
}

// NewConsulRegistry creates a new ConsulRegistry connected to the given Consul address.
func NewConsulRegistry(addr string) (*ConsulRegistry, error) {
	config := consul.DefaultConfig()
	config.Address = addr

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &ConsulRegistry{client: client}, nil
}

// Register registers a service instance with Consul, including an HTTP health check.
func (r *ConsulRegistry) Register(serviceID, serviceName, host string, port int) error {
	registration := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: host,
		Port:    port,
		Check: &consul.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", host, port),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	err := r.client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service %s: %w", serviceName, err)
	}

	log.Printf("Service %s registered with Consul (ID: %s, Address: %s:%d)", serviceName, serviceID, host, port)
	return nil
}

// Deregister removes a service instance from Consul.
func (r *ConsulRegistry) Deregister(serviceID string) error {
	err := r.client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service %s: %w", serviceID, err)
	}

	log.Printf("Service %s deregistered from Consul", serviceID)
	return nil
}

// Discover finds a healthy instance of a service by name and returns its address.
func (r *ConsulRegistry) Discover(serviceName string) (string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("failed to discover service %s: %w", serviceName, err)
	}

	if len(entries) == 0 {
		return "", fmt.Errorf("no healthy instances found for service %s", serviceName)
	}

	// Return the first healthy instance
	entry := entries[0]
	addr := fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port)
	return addr, nil
}
