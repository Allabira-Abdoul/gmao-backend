package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	httphandler "backend-gmao/apps/user-service/internal/adapters/primary/http"
	"backend-gmao/pkg/db"
	"backend-gmao/pkg/discovery"
	"github.com/gin-gonic/gin"
)

func main() {
	// Consul Config
	consulHost := os.Getenv("CONSUL_HOST")
	if consulHost == "" {
		consulHost = "127.0.0.1"
	}
	consulPort := os.Getenv("CONSUL_PORT")
	if consulPort == "" {
		consulPort = "8500"
	}
	consulURL := fmt.Sprintf("%s:%s", consulHost, consulPort)

	// Initialize Consul Registry
	registry, err := discovery.NewConsulRegistry(consulURL)
	if err != nil {
		log.Fatalf("Failed to create consul registry: %v", err)
	}

	// Service Config
	serviceID := fmt.Sprintf("user-service-%s", os.Getenv("INSTANCE_ID"))
	if os.Getenv("INSTANCE_ID") == "" {
		serviceID = "user-service-1"
	}
	serviceName := "user-service"
	host := os.Getenv("SERVICE_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	// Default port for local testing based on service
	port := 8080 // We will override this if PORT env var is set
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", &port)
	}

	// Database Connection (Placeholder)
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		_, err := db.InitPostgres(dsn)
		if err != nil {
			log.Printf("Warning: Failed to connect to DB: %v", err)
		} else {
			log.Println("Successfully connected to Postgres database")
		}
	}

	// Register with Consul
	err = registry.Register(serviceID, serviceName, host, port)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Register Health Route
	router.GET("/health", httphandler.HealthCheck)

	// Start Server
	go func() {
		addr := fmt.Sprintf(":%d", port)
		log.Printf("Starting %s on %s\n", serviceName, addr)
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down service...")

	if err := registry.Deregister(serviceID); err != nil {
		log.Printf("Failed to deregister service: %v", err)
	}
	log.Println("Service stopped")
}
