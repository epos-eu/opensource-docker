
# EPOS Open Source - Docker installer

## Introduction

EPOS Open Source - Docker installer is part of the EPOS Open Source project for local installation using Docker.
It contains a set of docker images to deploy the EPOS ecosystem. 

Use `epos-<os>-<architecture>` binary to spin up local environment on Linux, Mac OS X or Windows.

## Prerequisites

Docker Engine and Docker Compose installed on your host machine.
For further information follow the official guidelines: https://docs.docker.com/get-docker/

## Environment Variables
### Base environment configuration

| Name | Standard Value | Description |
|--|--|--|
| API_HOST | ${API_HOST} | API Host IP, if not set is generated automatically using machine IP |
| EXECUTE_HOST | ${API_HOST} | Internal variable to setup redirections for the external access service, if not set is generated automatically using machine IP |
| DEPLOY_PATH | / | Context path of the environment|
| BASE_CONTEXT | empty value | Context path name of the environment (similar to DEPLOY_PATH but without the initial /) |
| API_PATH | /api/v1 | API GATEWAY access path|
| GUI_PORT | 8000 | Port used by EPOS Data Portal or other GUIs |
| BACKOFFICE_GUI_PORT | 9000 | Port used by EPOS Backoffice UI or other Backoffice GUIs |
| API_PORT | 8080 | Port used by EPOS API Gateway |

### RabbitMQ configuration

| Name | Standard Value | Description |
|--|--|--|
| BROKER_USERNAME | changeme | RabbitMQ username |
| BROKER_PASSWORD | changeme | RabbitMQ password |
| BROKER_VHOST | changeme | RabbitMQ vhost |

### RabbitMQ configuration

| Name | Standard Value | Description |
|--|--|--|
| POSTGRES_USER | postgres | Database user |
| POSTGRESQL_PASSWORD | changeme | Database password |
| POSTGRES_DB | cerif | Database name |
| POSTGRESQL_CONNECTION_STRING | jdbc:postgresql://postgrescerif:5432/${POSTGRES_DB}?user=${POSTGRES_USER}&password=${POSTGRESQL_PASSWORD} | Database connection string based on previous configurations |
| PERSISTENCE_NAME | EPOSDataModel | Persistence Name of scientific metadata |
| PERSISTENCE_NAME_PROCESSING | EPOSProcessing | Persistence Name of processing metadata |

### Data Metadata Service configuration

| Name | Standard Value | Description |
|--|--|--|
| NUM_OF_PUBLISHERS | 10 | Number of publishers on rabbitmq |
| NUM_OF_CONSUMERS | 10 | Number of consumers on rabbitmq |
| CONNECTION_POOL_INIT_SIZE | 1 | Initial number of connections to database |
| CONNECTION_POOL_MIN_SIZE | 1 | Minimum number of connections to database |
| CONNECTION_POOL_MAX_SIZE | 20 | Maximum number of connections to database |

### Monitoring Service configuration

| Name | Standard Value | Description |
|--|--|--|
| MONITORING | false | True if activate interaction between system and monitoring service |
| MONITORING_URL | changeme | Monitoring service url |
| MONITORING_USER | changeme | Monitoring service username |
| MONITORING_PWD | changeme | Monitoring service password |

### Monitoring Service configuration

| Name | Standard Value | Description |
|--|--|--|
| DOCKER_REGISTRY | epos | Docker registry url |
| REGISTRY_USERNAME | changeme | Docker registry username |
| REGISTRY_PASSWORD | changeme | Docker registry password |

### GitLab/Hub configuration

| Name | Standard Value | Description |
|--|--|--|
| REPOTOKEN_DEFAULT | changeme | Token to access gitlab repositories, used from converter to download plugins |

### ICS-D configuration

| Name | Standard Value | Description |
|--|--|--|
| SWIRRL_BASE_PATH | https://epos-ics-d.brgm-rec.fr/swirrl-api/ | Temporary url to swirrl APIs |

### Docker Images for Open Source 

| Variable name | Image name | Default Tag |
|--|--|--|
| GUI_IMAGE | epos-gui | opensource |
| METADATA_DB_IMAGE | metadata-database-deploy | 2.2.0 |
| MESSAGE_BUS_IMAGE | rabbitmq | 3.11.7-management |
| GATEWAY_IMAGE | epos-api-gateway | 1.1.0 |
| RESOURCES_SERVICE_IMAGE | resources-service | 1.2.0 |
| DATA_METADATA_SERVICE_IMAGE | data-metadata-service | 2.3.5 |
| INGESTOR_IMAGE | ingestor-service | 1.3.0 |
| EXTERNAL_ACCESS_IMAGE | external-access-service | 1.2.0 |
| BACKOFFICE_SERVICE_IMAGE | backoffice-service | 2.1.0 |
| CONVERTER_IMAGE | converter-service | 1.1.5 |
| PROCESSING_ACCESS_SERVICE_IMAGE | distributed-processing-service | 0.2.0 |


## Maintenance

We regularly update images used in this stack.


## Installation

Download the binary file according to your OS.

Then give permissions on `epos-<os>-<architecture>` file from a Terminal (in Linux/MacOS):

```
chmod +x epos-<os>-<architecture>
```

## Usage

```
./epos-<os>-<architecture> <command>
```

The `<command>` field value is one of the following listed below:

```
EPOS Open Source CLI installer to deploy the EPOS System using docker-compose

Usage:
  [command] [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete an environment on docker
  deploy      Deploy an environment on docker
  export      Export configuration files in output folder, options: [env, compose]
  help        Help about any command
  populate    Populate the existing environment with metadata information

Flags:
  -h, --help   help for [command]

Use "[command] [command] --help" for more information about a command.
```

## Deploy a new environment

```
Deploy an enviroment with .env set up on docker

Usage:
  [command] deploy [flags]

Flags:
      --dockercompose string   --dockercompose 1
      --env string             --env 1
  -h, --help                   help for deploy
```

## Delete the existing environment

```
Delete an enviroment with .env set up on docker

Usage:
  [command] delete [flags]

Flags:
      --dockercompose string   --dockercompose 1
      --env string             --env 1
  -h, --help                   help for delete
```

## Populate the existing environment with metadata

Automatic option: download or create TTL files according to EPOS-DCAT-AP and use the following command:

```
Populate the existing environment with metadata information in a specific folder

Usage:
  [command] populate [flags]

Flags:
      --file string   file or folder
  -h, --help          help for populate
```

Manual option: use the API Gateway endpoint to manually ingest metadata TTL files into the catalogue.

## Export configuration file and docker-compose

```
Export configuration files for customization in output folder, options: [env, compose]

Usage:
  [command] export [flags]

Flags:
      --file string     File to export, available options: [env, compose]
  -h, --help            help for export
      --output string   Output folder
```


## Access URLs

EPOS Data Portal: 
```
http://<your-ip>:<GUI_PORT>/
```

EPOS Backoffice: 
```
http://<your-ip>:<BACKOFFICE_GUI_PORT>/
```

EPOS API Gateway: 
```
http://<your-ip>:<API_PORT>/
```
