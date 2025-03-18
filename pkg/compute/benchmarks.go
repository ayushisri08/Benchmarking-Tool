package compute

import (
	"fmt"
	"os/exec"
)

// GetBenchmarkResults retrieves benchmark results from a VM via SSH
func GetBenchmarkResults(projectID, zone, instanceName string) (string, error) {
	sshCommand := fmt.Sprintf(
		"gcloud compute ssh %s --project=%s --zone=%s --tunnel-through-iap --command 'cat /var/log/startupscript.log'",
		instanceName, projectID, zone,
	)

	fmt.Println("Retrieving Benchmark results...")
	output, err := exec.Command("bash", "-c", sshCommand).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run logs command: %s, error: %v", string(output), err)
	}

	return string(output), nil
}