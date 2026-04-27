package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
//
	// Service Config
	serviceID := fmt.Sprintf("api-gateway-%s", os.Getenv("INSTANCE_ID"))
	if os.Getenv("INSTANCE_ID") == "" {
		serviceID = "api-gateway-1"
	}
	serviceName := "api-gateway"
	host := os.Getenv("SERVICE_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := 8080
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", &port)
	}

	// Register the gateway itself with Consul
	err = registry.Register(serviceID, serviceName, host, port)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Gateway health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"service": "api-gateway",
		})
	})

	// Dynamic reverse proxy for all service routes
	// Pattern: /api/{service-name}/* -> {service-name}-service/
	// Example: GET /api/analytics/health -> analytics-service /health
	router.Any("/api/:service/*path", func(c *gin.Context) {
		serviceSuffix := c.Param("service")
		targetPath := c.Param("path")

		// Build the Consul service name
		targetServiceName := serviceSuffix + "-service"

		// Discover the service via Consul
		addr, err := registry.Discover(targetServiceName)
		if err != nil {
			log.Printf("Service discovery failed for %s: %v", targetServiceName, err)
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   "service_unavailable",
				"message": fmt.Sprintf("Could not discover service: %s", targetServiceName),
			})
			return
		}

		// Build target URL
		targetURL, err := url.Parse(fmt.Sprintf("http://%s", addr))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Failed to parse target URL",
			})
			return
		}

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Modify the request
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = targetPath
			req.URL.RawQuery = c.Request.URL.RawQuery
			req.Host = targetURL.Host

			// Forward original headers
			for key, values := range c.Request.Header {
				for _, value := range values {
					req.Header.Set(key, value)
				}
			}

			// Add gateway-specific headers
			req.Header.Set("X-Forwarded-For", c.ClientIP())
			req.Header.Set("X-Forwarded-Host", c.Request.Host)
			req.Header.Set("X-Gateway-Service", "api-gateway")
		}

		// Custom error handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error for %s: %v", targetServiceName, err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(fmt.Sprintf(`{"error":"proxy_error","message":"Failed to reach %s"}`, targetServiceName)))
		}

		log.Printf("Proxying request: %s %s -> %s%s", c.Request.Method, c.Request.URL.Path, addr, targetPath)
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// Also support /api/{service} without trailing path (e.g. /api/analytics)
	router.Any("/api/:service", func(c *gin.Context) {
		// Redirect to include a trailing slash
		c.Request.URL.Path = strings.TrimRight(c.Request.URL.Path, "/") + "/"
		router.HandleContext(c)
	})

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
	log.Println("Shutting down API Gateway...")

	if err := registry.Deregister(serviceID); err != nil {
		log.Printf("Failed to deregister service: %v", err)
	}
	log.Println("API Gateway stopped")
}
