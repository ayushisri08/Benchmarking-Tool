package compute

import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
)

// Service handles all compute operations
type Service struct {
	projectID   string
	zone        string
	sourceImage string
	networkName string
}

// NewService creates a new compute service
func NewService(projectID, zone, sourceImage, networkName string) *Service {
	return &Service{
		projectID:   projectID,
		zone:        zone,
		sourceImage: sourceImage,
		networkName: networkName,
	}
}

// ListInstances lists all VM instances in the specified project and zone
func (s *Service) ListInstances(w io.Writer) ([]string, error) {
	ctx := context.Background()
	
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create instances client: %w", err)
	}
	defer instancesClient.Close()
	
	req := &computepb.ListInstancesRequest{
		Project: s.projectID,
		Zone:    s.zone,
	}
	
	fmt.Fprintln(w, "Listing existing instances...")
	it := instancesClient.List(ctx, req)
	
	var instances []string
	for {
		instance, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list instances: %w", err)
		}
		
		fmt.Fprintf(w, " - Instance Name: %s\n", instance.GetName())
		instances = append(instances, instance.GetName())
	}
	
	if len(instances) == 0 {
		fmt.Fprintln(w, "No instances found.")
	}
	
	return instances, nil
}

// CreateInstance creates a new VM instance with the specified parameters
func (s *Service) CreateInstance(w io.Writer, instanceName, machineType string) error {
	ctx := context.Background()
	
	fmt.Fprintln(w, "Starting instance creation...")
	fmt.Fprintf(w, "Creating instance with name: %s in zone: %s\n", instanceName, s.zone)
	
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create instances client: %w", err)
	}
	defer instancesClient.Close()
	
	networkResourcePath := fmt.Sprintf("projects/%s/global/networks/%s", s.projectID, s.networkName)
	startupScript := getStartupScript()
	
	req := &computepb.InsertInstanceRequest{
		Project: s.projectID,
		Zone:    s.zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(10),
						SourceImage: proto.String(s.sourceImage),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", s.zone, machineType)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Network: proto.String(networkResourcePath),
					AccessConfigs: []*computepb.AccessConfig{
						{
							Name: proto.String("External NAT"),
							Type: proto.String("ONE_TO_ONE_NAT"),
						},
					},
				},
			},
			Metadata: &computepb.Metadata{
				Items: []*computepb.Items{
					{
						Key:   proto.String("startup-script"),
						Value: proto.String(startupScript),
					},
				},
			},
		},
	}
	
	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}
	
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("failed to wait for operation to complete: %w", err)
	}
	
	fmt.Fprintf(w, "Instance '%s' created successfully.\n", instanceName)
	return nil
}

// DeleteInstance deletes a VM instance
func (s *Service) DeleteInstance(w io.Writer, instanceName string) error {
	ctx := context.Background()
	
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create instances client: %w", err)
	}
	defer instancesClient.Close()
	
	req := &computepb.DeleteInstanceRequest{
		Project:  s.projectID,
		Zone:     s.zone,
		Instance: instanceName,
	}
	
	op, err := instancesClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete instance: %w", err)
	}
	
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for delete operation: %w", err)
	}
	
	fmt.Fprintf(w, "Instance '%s' deleted successfully.\n", instanceName)
	return nil
}

// getStartupScript returns the startup script for benchmarking
func getStartupScript() string {
	return `#!/bin/bash
set -x
exec > /var/log/startupscript.log 2>&1

# Update and install necessary packages
sudo apt-get update
sudo apt-get install -y sysbench stress-ng fio

# Run Sysbench CPU test
echo "Running Sysbench CPU test..."
sysbench cpu --cpu-max-prime=20000 run
echo "Sysbench CPU test completed!"

# Run Stress-ng
echo "Running Stress-ng..."
stress-ng --cpu 4 --timeout 60 --metrics-brief
echo "Stress-ng test completed!"

# Run the fio benchmarking tool
echo "Running fio disk benchmarks..."
fio --name=random-write --ioengine=posixaio --rw=randwrite --bs=4k --size=4g --numjobs=1 --runtime=60 --time_based --end_fsync=1
echo "Fio disk benchmark completed!"
`
}