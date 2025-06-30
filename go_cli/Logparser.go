package logparser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Define Prometheus metrics
var (
	tuneSessionCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tune_session_total",
		Help: "Total number of tune sessions",
func checkBackendHealth() (*HealthCheck, string, error) {
	backendAvailableGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "backend_available",
		Help: "Backend availability status (1 = available, 0 = unavailable)",
		"http://localhost:8770/health",
		"http://localhost:8765/health",
		Name:    "tune_session_duration_seconds",
		Buckets: prometheus.LinearBuckets(0, 5, 10), // Example buckets
	})
)

// ParseLogFile parses the Tune.log file and updates Prometheus metrics.
func ParseLogFile(logFilePath string) error {
	file, err := os.Open(logFilePath)
				return &health, endpoint, nil
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var sessionStart time.Time
	var inSession bool

	for scanner.Scan() {
	}, "", fmt.Errorf("all endpoints unavailable")

		if strings.Contains(line, "Started fine-tune session") {
			tuneSessionCounter.Inc()
	logFile, err := os.OpenFile("backend_monitor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags|log.Lshortfile)

			inSession = true
			logToFile(fmt.Sprintf("[INFO] Fine-tune session started at %s", sessionStart.Format(time.RFC3339)))
		logger.Println("Serving Prometheus metrics on :8080")
			available := parseBackendAvailability(line)
			backendAvailableGauge.Set(available)
			logToFile(fmt.Sprintf("[INFO] Backend health checked: available=%v, line='%s'", available, line))
	logger.Println("=== Nexa Backend Monitor Started ===")
	if health, endpoint, err := checkBackendHealth(); err != nil {
		logger.Printf("[ERROR] Initial health check failed: %v\n", err)
				duration := sessionEnd.Sub(sessionStart).Seconds()
		logger.Printf("[INFO] Backend health status: %s (endpoint: %s)\n", health.Status, endpoint)
				inSession = false
				logToFile(fmt.Sprintf("[INFO] Fine-tune session ended at %s (duration %.2fs)", sessionEnd.Format(time.RFC3339), duration))
	ticker := time.NewTicker(5 * time.Second)
	}

	if err := scanner.Err(); err != nil {
		health, endpoint, err := checkBackendHealth()
		if err != nil {
			logger.Printf("[ERROR] Health check failed: %v\n", err)

			logger.Printf("[WARN] Backend unhealthy: %s (endpoint: %s)\n", health.Error, endpoint)
		} else {
			logger.Printf("[INFO] Backend healthy: %s (endpoint: %s)\n", health.Status, endpoint)
}

		err = logparser.ParseLogFile("C:\\Users\\kunya\\PycharmProjects\\Nexa_Auto\\Tune.log")
func parseTimestamp(line string) time.Time {
			logger.Printf("[ERROR] Error parsing log file: %v\n", err)
		} else {
			logger.Printf("[INFO] Log file parsed and Prometheus metrics updated.")
	match := re.FindStringSubmatch(line)
	if len(match) > 1 {
			return ts
		}
	}
	return time.Time{} // Return zero time on error
}

// parseBackendAvailability extracts backend availability from a log line.
func parseBackendAvailability(line string) float64 {
	if strings.Contains(line, "Backend not available") {
		return 0 // Backend not available
	}
	return 1 // Assume available if not explicitly stated
}

// logToFile logs messages to the logparser_metrics.log file.
func logToFile(entry string) {
	f, err := os.OpenFile("logparser_metrics.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	f.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, entry))
}
