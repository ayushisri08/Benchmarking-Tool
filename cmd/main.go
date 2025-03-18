// package main

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"benchmarking/pkg/compute"
// 	"benchmarking/pkg/config"
// )

// func main() {
// 	// Load default configuration
// 	cfg := config.DefaultConfig()

// 	// Check if project ID is set
// 	if cfg.ProjectID == "" {
// 		fmt.Print("Enter your GCP Project ID: ")
// 		fmt.Scan(&cfg.ProjectID)
// 	}

// 	// List existing instances
// 	instances, err := compute.ListInstances(os.Stdout, cfg.ProjectID, cfg.Zone)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error listing instances: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// Get user action
// 	var action string
// 	fmt.Println("Enter 'create' to create an instance or 'delete' to delete an instance:")
// 	fmt.Scan(&action)

// 	switch action {
// 	case "create":
// 		// Get instance details
// 		fmt.Println("Enter the desired name for the new instance:")
// 		var instanceName string
// 		fmt.Scan(&instanceName)

// 		fmt.Println("Enter the type of instance you want to create (e.g., e2-micro):")
// 		var machineType string
// 		fmt.Scan(&machineType)

// 		// Create instance
// 		if err := compute.CreateInstance(os.Stdout, cfg.ProjectID, cfg.Zone, instanceName, machineType, cfg.SourceImage, cfg.NetworkName); err != nil {
// 			fmt.Fprintf(os.Stderr, "Error creating instance: %v\n", err)
// 			os.Exit(1)
// 		}

// 		// Wait for benchmarks to run
// 		fmt.Println("Waiting for benchmarks to complete (120 seconds)...")
// 		time.Sleep(120 * time.Second)

// 		// Get benchmark results
// 		results, err := compute.GetBenchmarkResults(cfg.ProjectID, cfg.Zone, instanceName)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "Error retrieving benchmark results: %v\n", err)
// 		} else {
// 			fmt.Println("Benchmark Results:")
// 			fmt.Println(results)
// 		}

// 	case "delete":
// 		if len(instances) == 0 {
// 			fmt.Println("No instances available to delete.")
// 			os.Exit(0)
// 		}

// 		fmt.Println("Enter the name of the instance to delete from the following list:")
// 		for _, name := range instances {
// 			fmt.Printf(" - %s\n", name)
// 		}

// 		var instanceName string
// 		fmt.Scan(&instanceName)
// 		if err := compute.DeleteInstance(os.Stdout, cfg.ProjectID, cfg.Zone, instanceName); err != nil {
// 			fmt.Fprintf(os.Stderr, "Error deleting instance: %v\n", err)
// 		}

// 	default:
// 		fmt.Println("Invalid action. Please specify 'create' or 'delete'.")
// 	}
// }

package main

import (
	"fmt"
	"os"
	"time"

	"benchmarking/pkg/compute"
	"benchmarking/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.DefaultConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		fmt.Print("Enter your GCP Project ID: ")
		fmt.Scan(&cfg.ProjectID)
	}

	// List existing instances
	instances, err := compute.ListInstances(os.Stdout, cfg.ProjectID, cfg.Zone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing instances: %v\n", err)
		os.Exit(1)
	}

	// Get user action
	var action string
	fmt.Println("Enter 'create' to create an instance or 'delete' to delete an instance:")
	fmt.Scan(&action)

	switch action {
	case "create":
		// Get instance details
		fmt.Println("Enter the desired name for the new instance:")
		var instanceName string
		fmt.Scan(&instanceName)

		fmt.Println("Enter the type of instance you want to create (e.g., e2-micro):")
		var machineType string
		fmt.Scan(&machineType)

		// Create instance
		if err := compute.CreateInstance(os.Stdout, cfg.ProjectID, cfg.Zone, instanceName, machineType, cfg.SourceImage, cfg.NetworkName); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating instance: %v\n", err)
			os.Exit(1)
		}

		// Wait for benchmarks to run
		fmt.Printf("Waiting for benchmarks to complete (%d seconds)...\n", int(cfg.WaitTime.Seconds()))
		select {
		case <-time.After(cfg.WaitTime):
			// Continue after timeout
		}

		// Get benchmark results
		results, err := compute.GetBenchmarkResults(cfg.ProjectID, cfg.Zone, instanceName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving benchmark results: %v\n", err)
		} else {
			fmt.Println("Benchmark Results:")
			fmt.Println(results)
		}

	case "delete":
		if len(instances) == 0 {
			fmt.Println("No instances available to delete.")
			os.Exit(0)
		}

		fmt.Println("Enter the name of the instance to delete from the following list:")
		for _, name := range instances {
			fmt.Printf(" - %s\n", name)
		}

		var instanceName string
		fmt.Scan(&instanceName)
		if err := compute.DeleteInstance(os.Stdout, cfg.ProjectID, cfg.Zone, instanceName); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting instance: %v\n", err)
		}

	default:
		fmt.Println("Invalid action. Please specify 'create' or 'delete'.")
	}
}