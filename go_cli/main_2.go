package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"C:\\Users\\kunya\\PycharmProjects\\Nexa_Auto\\go_cli\\logparser"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HealthCheck struct {
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
	Timestamp string `json:"timestamp"`
}

func checkBackendHealth() (*HealthCheck, error) {
	endpoints := []string{
		"http://localhost:8000/health",
		"http://127.0.0.1:8000/health",
		"http://localhost:8770/health", // trainer server
		"http://localhost:8765/health", // session server
	}

	for _, endpoint := range endpoints {
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get(endpoint)
		if err == nil && resp.StatusCode == 200 {
			var health HealthCheck
			if err := json.NewDecoder(resp.Body).Decode(&health); err == nil {
				resp.Body.Close()
				return &health, nil
			}
			resp.Body.Close()
		}
	}

	return &HealthCheck{
		Status:    "error",
		Error:     "All backend endpoints unavailable",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

func main() {
	// Start the Prometheus HTTP server.
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Serving Prometheus metrics on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Initial health check
	if health, err := checkBackendHealth(); err != nil {
		log.Printf("Initial health check failed: %v\n", err)
	} else {
		log.Printf("Backend health status: %s\n", health.Status)
	}

	// Parse the log file periodically
	ticker := time.NewTicker(5 * time.Second) // Adjust interval as needed
	defer ticker.Stop()

	for range ticker.C {
		if health, err := checkBackendHealth(); err != nil {
			log.Printf("Health check failed: %v\n", err)
		} else if health.Status != "ok" {
			log.Printf("Backend unhealthy: %s\n", health.Error)
		}

		err := logparser.ParseLogFile("C:\\Users\\kunya\\PycharmProjects\\Nexa_Auto\\Tune.log")
		if err != nil {
			log.Printf("Error parsing log file: %v\n", err)
		}
	}

	select {} // Keep the application running
}
