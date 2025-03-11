package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yourusername/vm-manager/pkg/compute"
	"github.com/yourusername/vm-manager/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Create compute service
	computeService := compute.NewService(
		cfg.ProjectID,
		cfg.Zone,
		cfg.SourceImage,
		cfg.NetworkName,
	)

	reader := bufio.NewReader(os.Stdin)

	// List existing instances
	instances, err := computeService.ListInstances(os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing instances: %v\n", err)
		os.Exit(1)
	}

	// Get action from user
	fmt.Println("\nEnter an action (create, delete, list):")
	action, _ := reader.ReadString('\n')
	action = strings.TrimSpace(action)

	if !config.IsValidAction(action) {
		fmt.Println("Invalid action. Please specify 'create', 'delete', or 'list'.")
		os.Exit(1)
	}

	switch strings.ToLower(action) {
	case "list":
		// Already listed instances above
		return

	case "create":
		// Get instance name
		fmt.Println("Enter the desired name for the new instance:")
		instanceName, _ := reader.ReadString('\n')
		instanceName = strings.TrimSpace(instanceName)

		// Get machine type
		fmt.Println("Enter the type of instance you want to create (e.g., e2-micro):")
		machineType, _ := reader.ReadString('\n')
		machineType = strings.TrimSpace(machineType)

		// Create the instance
		if err := computeService.CreateInstance(os.Stdout, instanceName, machineType); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating instance: %v\n", err)
			os.Exit(1)
		}

		// Wait for the instance to potentially complete benchmarks
		fmt.Println("Waiting for benchmarks to complete (120 seconds)...")
		time.Sleep(120 * time.Second)

		// Retrieve and print benchmark results
		results, err := computeService.RunBenchmarkResults(instanceName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving benchmark results: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Benchmark Results:")
		fmt.Println(results)

	case "delete":
		if len(instances) == 0 {
			fmt.Println("No instances available to delete.")
			os.Exit(0)
		}

		fmt.Println("Enter the name of the instance to delete from the list above:")
		instanceName, _ := reader.ReadString('\n')
		instanceName = strings.TrimSpace(instanceName)

		if err := computeService.DeleteInstance(os.Stdout, instanceName); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting instance: %v\n", err)
			os.Exit(1)
		}
	}
}