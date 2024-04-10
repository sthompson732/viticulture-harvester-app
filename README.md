# Viticulture Data Harvester Application

## Overview

The Viticulture Harvester is a proof of concept application that demonstrates foundational functionalities for vineyard management using Go. It provides a robust platform with RESTful endpoints for CRUD operations on vineyard data, integrated with PostgreSQL and cloud storage solutions, and utilizes Google Cloud Scheduler for regular data ingestion tasks.

## Technical Configuration

- **Language**: Go 1.21
- **Database**: PostgreSQL, for robust data management and querying capabilities.
- **Storage**: Integrated with cloud storage solutions to manage satellite images and other large datasets.
- **Scheduling**: Uses Google Cloud Scheduler for automating data ingestion processes.
- **Architecture**: RESTful API with foundational middleware for logging, authentication, and authorization.
- **Dependencies**: Managed with Go Modules.
- **Configuration**: Settings and configurations managed via `config.yaml`.
- **Main file location**: `cmd/harvester/main.go`

## Features

- **Vineyard Management**: Supports CRUD operations to create, update, delete, and retrieve vineyard records efficiently.
- **Image Processing**: Framework to handle satellite image data processing and storage.
- **Data Reporting**: Capabilities to generate comprehensive reports summarizing vineyard status and health.
- **Scheduled Tasks**: Automated data ingestion and processing tasks configured through Google Cloud Scheduler.

## Getting Started

### Prerequisites

Before you begin, ensure you have installed the following software on your system:

- [Git](https://git-scm.com/downloads) (latest version)
- [Go](https://go.dev/dl/) (version 1.21 or higher)
- An IDE of your choice (VS Code, GoLand, etc.)

### Project Structure

```text
/viticulture-harvester-app
    /cmd
        /harvester                  # Main application entry point
            main.go
    /internal                       # Application's internal logic
        /api                        # API endpoints and server setup using Gin
            router.go               # Gin router setup
            handlers.go             # API handlers
            router_test.go          # Tests for the Gin router setup
            handlers_test.go        # Tests for the API handlers
        /clients                    # External service clients
            satellite.go            # Client for satellite data API
            satellite_test.go       # Tests for satellite client
            soil.go                 # Client for soil data API
            soil_test.go            # Tests for soil client
        /config                     # Configuration management
            config.go               # Centralized configuration logic
            config_test.go          # Tests for configuration logic
        /db                         # Database Interactions
            db.go                   # Postgres Interaction
            db_test.go              # Tests for database logic
        /model                      # Data models
            models.go               # Structs for database and JSON
            models_test.go          # Tests for data models
        /service                    # Service Interfaces
            imageservice.go         # Service for Satellite Image Interface
            soildataservice.go      # Service for Soil Data Interface
            vineyardservice.go      # Service for Vineyard Data Interface
        /storage                    # Cloud Storage Interactions
            storage.go              # GCP Cloud Storage Interaction
            storage_test.go         # Tests for storage logic
    /pkg                            # External packages that might be used across projects
        /util                       # Utility functions
            util.go                 # Helper functions and utilities
            util_test.go            # Tests for utility functions
    /configs                        # Configuration files
        config.yaml                 # Application configurations
    /scripts                        # Auxiliary scripts, e.g., for setup
        setup.sh                    # Setup script for initializing the application
    go.mod                          # Go module definitions
    go.sum                          # Go module checksums
    Dockerfile                      # Dockerfile for containerizing the application
    .gitignore                      # Git ignore file
    README.md                       # Project overview and setup instructions
```

### Environment Setup

Follow these steps to set up the project:

#### 1. Clone the Repository

Clone the repository to your local machine.

```
git clone https://github.com/sthompson732/viticulture-harvester-app.git
cd viticulture-harvester-app
```

#### 2. Install Go

Download and install Go (if not already installed). Visit the [official Go download page](https://go.dev/dl/) and follow the installation instructions for your operating system.

#### 3. Set Up Go Environment

Configure your Go environment by setting the `GOPATH` and updating the `PATH` variable:

For Windows:

```
setx GOPATH "%USERPROFILE%\go"
setx PATH "%PATH%;%USERPROFILE%\go\bin;%GOROOT%\bin"
```

For Linux or macOS:

```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin:$GOROOT/bin
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export PATH=$PATH:$GOPATH/bin:$GOROOT/bin" >> ~/.bashrc
source ~/.bashrc
```

#### 4. Navigate to the Project Directory

Change to the directory where the main application is located:

```
cd cmd/harvester
```

#### 5. Install Dependencies

Fetch and install the project dependencies:

```
go mod tidy
```

This command cleans up the module's dependencies, adding missing ones and removing unused ones.

#### 6. Build the Application

Compile the application into an executable:

```
go build -o harvester.exe  # On Windows
go build -o harvester      # On Linux or macOS
```

#### 7. Run the Application

Run the compiled application:

```
./harvester.exe  # On Windows
./harvester      # On Linux or macOS
```

### Next Steps

- Explore the application features.
- Check the [issues tab](<link-to-issues-tab>) on GitHub for current issues and feature requests.

### Support

If you encounter any problems or have suggestions, please file an issue on the GitHub repository.