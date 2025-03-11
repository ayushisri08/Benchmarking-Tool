package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	ProjectID   string
	Zone        string
	SourceImage string
	NetworkName string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("PROJECT_ID environment variable is required")
	}

	zone := os.Getenv("ZONE")
	if zone == "" {
		zone = "us-central1-a" // Default zone
	}

	sourceImage := os.Getenv("SOURCE_IMAGE")
	if sourceImage == "" {
		sourceImage = "projects/debian-cloud/global/images/family/debian-11"
	}

	networkName := os.Getenv("NETWORK_NAME")
	if networkName == "" {
		networkName = "default-network"
	}

	return &Config{
		ProjectID:   projectID,
		Zone:        zone,
		SourceImage: sourceImage,
		NetworkName: networkName,
	}, nil
}

// GetNetworkPath returns the full path to the network resource
func (c *Config) GetNetworkPath() string {
	return fmt.Sprintf("projects/%s/global/networks/%s", c.ProjectID, c.NetworkName)
}

// GetMachineTypePath returns the full path to the machine type
func (c *Config) GetMachineTypePath(machineType string) string {
	return fmt.Sprintf("zones/%s/machineTypes/%s", c.Zone, machineType)
}

// IsValidAction checks if the provided action is valid
func IsValidAction(action string) bool {
	action = strings.ToLower(action)
	return action == "create" || action == "delete" || action == "list"
}