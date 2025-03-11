package compute

import (
	"fmt"
	"os/exec"
)

// RunBenchmarkResults retrieves benchmark results from a VM instance via SSH
func (s *Service) RunBenchmarkResults(instanceName string) (string, error) {
	// Command to retrieve benchmark results via SSH
	sshCommand := fmt.Sprintf(
		"gcloud compute ssh %s --project=%s --zone=%s --tunnel-through-iap --command 'cat /var/log/startupscript.log'",
		instanceName, s.projectID, s.zone,
	)
	
	fmt.Println("Retrieving Benchmark results...")
	
	output, err := exec.Command("bash", "-c", sshCommand).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run logs command: %s, error: %w", string(output), err)
	}
	
	return string(output), nil
}