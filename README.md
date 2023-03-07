
# EPOS Open Source - Docker installer

## Introduction

EPOS Open Source - Docker installer is part of the EPOS Open Source project for local installation using Docker.
It contains a set of docker images to deploy the EPOS ecosystem. 

Use `epos-open-source` binary to spin up local environment on Linux, Mac OS X or Windows.

## Prerequisites

Docker Engine and Docker Compose installed on your host machine.
For further information follow the official guidelines: https://docs.docker.com/get-docker/

## Stack

The EPOS Open Source for local installation with docker stack consist of the following containers:

| Container       | Versions           | Image                              | 
|-----------------|--------------------|------------------------------------|
| EPOS Data Portal       | 1.0.3           | epos-ci.brgm.fr:5005/epos/epos-gui/deploy:1-0-3                              | 
| EPOS Backoffice       | x.y.z           | epos-ci.brgm.fr:5005/epos/epos-gui/deploy:x.y.z                             | 
| API Gateway       | 1-0-12           | epos-ci.brgm.fr:5005/epos/simple-gateway:1-0-13                              | 
| Resources Service       | 1-1-3           | epos-ci.brgm.fr:5005/epos/resources-service:1-1-3                              | 
| External Access Service       | 1-0-11           | epos-ci.brgm.fr:5005/epos/external-access-service:1-0-11external-access-service:1-0-11                              | 
| Distributed Processing Service       | 0-0-5           | epos-ci.brgm.fr:5005/epos/distributed-processing-service:0-0-5                              | 
| Ingestor Service       | 1-0-7           | epos-ci.brgm.fr:5005/epos/ingestor-service:1-0-7 
| Backoffice Service       | 1-0-3           | epos-ci.brgm.fr:5005/epos/backoffice-service:1-0-3                              | 
| Converter Service       | 1-0-19           | epos-ci.brgm.fr:5005/epos/converter-service:1-0-19                                    |                        | 
| Data Metadata Service       | 1-2-20           | epos-ci.brgm.fr:5005/epos/data-metadata-service:1-2-20                              | 
| Metadata Cache      | 0-21-2           | epos-ci.brgm.fr:5005/epos/metadata-cache:0-21-2                            | 
| Rabbitmq      | 3.10.7-management           | rabbitmq:3.10.7-management                             | 


## Maintenance

We regularly update images used in this stack and release them together.


## Installation

Clone this repository or download it.

Copy the file `/template-configuration/configuration.sh` into `/configuration` folder, the update the environment variables declared in that file with your own setup.

Pay attenction to the GITLAB_* variables:

```
# GITLAB LOGIN
export GITLAB_USERNAME="changeme"
export GITLAB_PASSWORD="changeme"

```

Then give permissions on `epos-open-source` file from a Terminal in Linux/MacOS:

```
chmod +x epos-open-source
```

## Usage

```
./epos-open-source <argument>
```

The `<argument>` field value is one of the following listed below:

```
deploy: deploy the environment

remove: remove the environment

populate: run ingestion phase

restartconverter: restart converter service to sync metadata with plugins

showhost: print hostpath

install: install docker and docker-compose, it requires a system restart
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

## Example

Deploy a new enviroment:

```
./epos-open-source deploy
```

Ingest metadata information from metadata cache:

```
./epos-open-source populate
```

Remove an existing environment:

```
./epos-open-source remove
```

## Next updates

- [ ] Ingest a single file
- [ ] Restart a single component
- [ ] Easy logs from cli
