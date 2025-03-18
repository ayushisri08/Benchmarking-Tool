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

// ListInstances lists all VM instances in the specified project and zone
func ListInstances(w io.Writer, projectID, zone string) ([]string, error) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create instances client: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.ListInstancesRequest{
		Project: projectID,
		Zone:    zone,
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

// CreateInstance creates a new VM instance with benchmark tools
func CreateInstance(w io.Writer, projectID, zone, instanceName, machineType, sourceImage, networkName string) error {
	ctx := context.Background()
	fmt.Fprintln(w, "Starting instance creation...")
	fmt.Fprintf(w, "Creating instance with name: %s in zone: %s\n", instanceName, zone)

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create instances client: %w", err)
	}
	defer instancesClient.Close()

	networkResourcePath := fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName)
	startupScript := getStartupScript()

	req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(10),
						SourceImage: proto.String(sourceImage),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)),
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
func DeleteInstance(w io.Writer, projectID, zone, instanceName string) error {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.DeleteInstanceRequest{
		Project:  projectID,
		Zone:     zone,
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

// getStartupScript returns the script for benchmarking
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
echo "Running Fio disk test..."
fio --name=random-write --ioengine=posixaio --rw=randwrite --bs=4k --size=4g --numjobs=1 --runtime=60 --time_based --end_fsync=1
echo "Fio disk test completed!"
`
}