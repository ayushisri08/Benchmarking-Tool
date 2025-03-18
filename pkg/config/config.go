
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// GCPConfig holds configuration for GCP interactions
type GCPConfig struct {
	ProjectID   string
	Zone        string
	NetworkName string
	SourceImage string
	WaitTime    time.Duration
}

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig() (*GCPConfig, error) {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Warning: .env file not found or cannot be loaded: %v\n", err)
	}

	// Get wait time from env or use default
	waitTimeStr := getEnv("GCP_WAIT_TIME", "120")
	waitTimeSec, err := strconv.Atoi(waitTimeStr)
	if err != nil {
		fmt.Printf("Warning: invalid GCP_WAIT_TIME value, using default: %v\n", err)
		waitTimeSec = 120
	}

	config := &GCPConfig{
		ProjectID:   getEnv("GCP_PROJECT_ID", ""),
		Zone:        getEnv("GCP_ZONE", "us-central1-a"),
		NetworkName: getEnv("GCP_NETWORK_NAME", "default-network"),
		SourceImage: getEnv("GCP_SOURCE_IMAGE", "projects/debian-cloud/global/images/family/debian-11"),
		WaitTime:    time.Duration(waitTimeSec) * time.Second,
	}

	return config, nil
}

// DefaultConfig returns default configuration, checking environment variables
func DefaultConfig() *GCPConfig {
	config, _ := LoadConfig()
	return config
}

// Validate checks if the configuration is valid
func (c *GCPConfig) Validate() error {
	if strings.TrimSpace(c.ProjectID) == "" {
		return fmt.Errorf("project ID is required")
	}
	return nil
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}