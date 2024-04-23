# Viticulture Data Harvester Application

## Overview

The Viticulture Harvester provides a comprehensive toolset designed for advanced vineyard management. Utilizing Go, PostgreSQL, and the Google Cloud Platform, this application enhances operational decision-making through detailed environmental data and high-resolution imagery integration.

### Purpose and Design:
Created to support precision agriculture, this platform automates and synthesizes data from multiple sources including satellite imagery and environmental sensors. This integration offers vineyard managers a dynamic view of their fields, facilitating precise cultivation strategies and resource management.

### Infrastructure and Technology:
Employing robust Google Cloud services and PostgreSQL, the Viticulture Harvester scales seamlessly to accommodate extensive data sets with optimal performance. Regular updates managed via Google Cloud Scheduler ensure timely data refreshes, making critical information readily available through efficient RESTful APIs.

This system not only simplifies comprehensive vineyard management but also equips operators with the tools to significantly enhance productivity and sustainability.

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

- **Vineyard Management**: Provides comprehensive tools for managing vineyard records, enabling CRUD operations (create, read, update, delete) that ensure detailed tracking and data accuracy.
- **High-Resolution Imagery Handling**: Supports the processing and storage of satellite and orthophoto imagery to aid in detailed monitoring and analysis of vineyard conditions.
- **Data Reporting**: Offers capabilities to generate detailed reports, summarizing vineyard health and activities, which supports informed decision-making and strategic planning.
- **Automated Data Ingestion**: Uses Google Cloud Scheduler to automate the regular collection of imagery and sensor data, keeping vineyard information consistently up-to-date.
- **IoT Integration**: Incorporates data from IoT sensors to provide real-time insights into vineyard conditions, optimizing resource management and operational responses.
- **API-Driven Interactions**: Facilitates robust API endpoints for efficient data retrieval and integration, enabling seamless interactions with external systems and applications.

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
    /harvester
        main.go                # Initializes services and starts the server.
/configs
    config.yaml                # Contains all application configurations.
/internal
    /api
        router.go              # Sets up HTTP routes and connects them with handlers.
        handlers.go            # Processes requests and returns responses.
    /clients
        satelliteclient.go      # Handles requests to satellite data APIs.
        soilclient.go           # Handles requests to soil data APIs.
        pestclient.go           # Handles requests to pest data APIs.
        weatherclient.go        # Handles requests to weather data APIs.
    /config
        config.go              # Loads and parses the config.yaml file.
    /db
        db.go                  # Manages database interactions.
    /model
        models.go              # Structures corresponding to database tables.
    /scheduler
        scheduler.go           # Manages timed data fetching jobs.
    /server
        server.go              # Configures and runs the HTTP server.
    /service
        imageservice.go        # Manages image data operations.
        pestservice.go         # Manages pest data operations.
        satelliteservice.go    # Manages satellite imagery operations.
        soilservice.go         # Manages soil data operations.
        vineyardservice.go     # Manages vineyard data operations.
        weatherservice.go      # Manages weather data operations.
    /storage
        storage.go             # Manages file storage operations.
/pkg
    /util
        util.go                # Provides common utility functions.
/scripts
    setup.sh                   # Sets up the application environment.
    /sql
        initdb.sql             # Initializes the database with the required schema.
        seed_data.sql          # Populates the database with initial data.
Dockerfile                     # For building the application's Docker container.
.gitignore                     # Specifies files to ignore in git operations.
go.mod                         # Manages dependencies.
go.sum                         # Contains hashes of dependencies for verification.
README.md                      # Project documentation and setup instructions.

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

#### 6. Database Setup

To initialize the database, run the following command from the project root:

``` bash
psql -U username -d databasename -a -f ./scripts/sql/initdb.sql

```

Replace username and databasename with your PostgreSQL username and database name, respectively.

##### Seed Data

If you need to seed your database with initial data, run:

``` bash
psql -U username -d databasename -a -f ./scripts/sql/seeddata.sql
```

#### 7. Build the Application

Compile the application into an executable:

```
go build -o harvester.exe  # On Windows
go build -o harvester      # On Linux or macOS
```

#### 8. Run the Application

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