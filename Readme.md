# Google Cloud VM Manager

A command-line tool for managing Google Cloud VM instances with built-in benchmarking capabilities.

## Features

- List existing VM instances in your Google Cloud project
- Create new VM instances with custom specifications
- Run benchmarks automatically on instance creation (CPU, stress tests, disk I/O)
- View benchmark results directly from the console
- Delete instances when no longer needed

## Prerequisites

- Go 1.19 or higher
- Google Cloud SDK installed and configured
- GCP project with Compute Engine API enabled
- Proper GCP credentials set up

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/yourusername/vm-manager.git
   cd vm-manager
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Create a `.env` file with your configuration (or set environment variables):
   ```
   PROJECT_ID=your-project-id
   ZONE=us-central1-a
   SOURCE_IMAGE=projects/debian-cloud/global/images/family/debian-11
   NETWORK_NAME=default-network
   ```

4. Build the application:
   ```
   go build -o vm-manager
   ```

## Usage

Run the application:
```
./vm-manager
```

Follow the interactive prompts to:
1. View existing instances
2. Create a new instance with benchmarking
3. Delete an existing instance

## Benchmarks Performed

When creating a new instance, the following benchmarks are automatically run:
- Sysbench CPU test
- Stress-ng CPU load test
- Fio disk I/O benchmark