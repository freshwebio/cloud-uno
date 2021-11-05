![Cloud Uno](/resources/logo.svg)

Cloud::1 enables running cloud services on your local machine for faster, cheaper, end-to-end development.

*If you were wondering, the name Cloud::1 (cloud uno) is derived from the ::1 IPv6 address that is the shorthand for the loopback localhost ip.*

## Prerequisites

### Running as a container or standalone

- [Docker](https://docs.docker.com/get-docker/) - Docker is required to run Cloud::1 as a container and to orchestrate services that run in docker containers.

### Build yourself

- [Docker](https://docs.docker.com/get-docker/) - Some services in Cloud::1 run in their own docker containers so even when you run the app as a standalone binary, you still need Docker.
- [Go](https://golang.org/dl/) - Cloud::1 is implemented in Go so in order to build from source you need Go installed.

## Configuration

The following provides configuration for all Cloud::1 applications.

All configuration should be provided via environment variables, command line options or a configuration file.

### File System

**(optional)**

 The file system to use for Cloud::1 custom emulators, the only choice is `memory`. Any other value will mean the os is used.

**Type** string

| Source          | Example                        |
| --------------- | :----------------------------- |
| **Flag**        | -cloud_uno_file_system memory  |
| **Environment** | CLOUD_UNO_FILE_SYSTEM=memory   |
| **File**        | cloud_uno_file_system memory   |

### Data Directory

**(optional, default = `/lib/data`)**

The location in the docker container to store all the data for the cloud service emulators. (If your file system is set to memory, then not all data for services will be persisted to disk)

**Type** string

| Source          | Example                            |
| --------------- | :--------------------------------- |
| **Flag**        | -cloud_uno_data_dir /path/to/data  |
| **Environment** | CLOUD_UNO_DATA_DIR=/path/to/data   |
| **File**        | cloud_uno_data_dir /path/to/data   |

### Run on Host

**(optional)**

If set, this will enable in-process features that are only available to privileged host applications. It will also embed the functionality to interact with the os hosts file in-process as opposed to in docker mode where a host agent that communicates over a unix socket is needed.

**Type** bool

| Source          | Example                      |
| --------------- | :--------------------------- |
| **Flag**        | -cloud_uno_run_on_host true  |
| **Environment** | CLOUD_UNO_RUN_ON_HOST=true   |
| **File**        | cloud_uno_run_on_host true   |

### Server IP

**(optional, default = `172.18.0.22`)**

The IP Address the cloud uno server is running on, this is ignored when running the server directly on the host.

**It is down to you to make sure the server sits behind the configured IP!**

**Type** string

| Source          | Example                      |
| --------------- | :--------------------------- |
| **Flag**        | -cloud_uno_ip 172.18.0.24    |
| **Environment** | CLOUD_UNO_IP=172.18.0.24     |
| **File**        | cloud_UNO_ip 172.18.0.24     |

### Hosts File Path

**(optional)**

A custom path to the hosts file on the host machine, otherwise defaults to the correct hosts file for the OS the host agent/server directly on the host is running on.

**Type** string

| Source          | Example                              |
| --------------- | :----------------------------------- |
| **Flag**        | -cloud_uno_hosts_path /custom/hosts  |
| **Environment** | CLOUD_UNO_HOSTS_PATH=/custom/hosts   |
| **File**        | cloud_uno_hosts_path /custom/hosts   |

### AWS Services

**(required if Google Cloud and Azure services aren't provided)**

AWS Services to run emulations for.

**Type** string

| Source          | Example                              |
| --------------- | :----------------------------------- |
| **Flag**        | -cloud_uno_aws_services s3,dynamodb  |
| **Environment** | CLOUD_UNO_AWS_SERVICES=s3,dynamodb   |
| **File**        | cloud_uno_aws_services s3,dynamodb   |

### Google Cloud Services

**(required if AWS and Azure services aren't provided)**

Google Cloud Services to run emulations for.

**Type** string

| Source          | Example                                            |
| --------------- | :------------------------------------------------- |
| **Flag**        | -cloud_uno_gcloud_services cloudstorage,datastore  |
| **Environment** | CLOUD_UNO_GCLOUD_SERVICES=cloudstorage,datastore   |
| **File**        | cloud_uno_gcloud_services cloudstorage,datastore   |

### Azure Services

**(required if AWS and Google Cloud services aren't provided)**

Azure Services to run emulations for.

**Type** string

| Source          | Example                                   |
| --------------- | :---------------------------------------- |
| **Flag**        | -cloud_uno_azure_services storage,cosmos  |
| **Environment** | CLOUD_UNO_AZURE_SERVICES=storage,cosmos   |
| **File**        | cloud_uno_azure_services storage,cosmos   |

### Debug

**(optional)**

Whether or not to run the application in debug mode, where debug logs will be written to stdout.

**Type** bool

| Source          | Example     |
| --------------- | :---------- |
| **Flag**        | -debug true |
| **Environment** | DEBUG=true  |
| **File**        | debug true  |

## Setting up

### Running In Docker

TODO: Provide instructions for running as a Docker container and with Docker compose.

When running Cloud::1 in Docker, in order to use networking/load balancing features of the cloud platform emulators you will need to configure a network and assign a static ip to your instance of Cloud::1. 

The default static ip that is assumed to be used for Cloud::1 is `172.18.0.22` but you can configure it to be whatever you like, you just need to make
sure the host agent is configured with the correct static ip.

Docker compose example:

```yaml
version: "3.9"
services:
  clouduno:
    image: clouduno
    environment:
      # Set to use an in-memory file system
      # for all cloud::1 services that are custom emulators.
      # (This won't take effect for most emulator APIs as they are wrappers around open source software or vendor-provided emulators)
      CLOUD_UNO_FILE_SYSTEM: "memory" # Defaults to OS
      CLOUD_UNO_DATA_DIR: /lib/path/to/custom/data # Defaults to /lib/data
      # Enable to show all the debug logs.
      DEBUG: true
      # AWS services to run.
      CLOUD_UNO_AWS_SERVICES: s3,dynamodb
      # Google cloud services to run.
      CLOUD_UNO_GCLOUD_SERVICES: cloudstorage,datastore
      # Azure services to run.
      CLOUD_UNO_AZURE_SERVICES: storage,cosmos
    ports:
      # Expose on port 80 on the host as due to the static ip
      # there won't be any conflicts.
      - "80:5988"
    networks:
      clouduno:
        ipv4_address: 172.18.0.22
    volumes:
     - 'host/path/to/custom/data:/lib/path/to/custom/data'
      - '/var/run/docker.sock:/var/run/docker.sock'
      # In order to make use of this you will need to make sure the host agent
      # has been installed and is running.
      - '/var/run/clouduno.sock:/var/run/clouduno.sock'
networks:
  clouduno:
    driver: bridge
    ipam:
      driver: default
      config:
        - subnet: 172.18.0.0/16
```

### Host Agent

When running Cloud::1 in Docker you need to run a host agent that deals with updating your machine's host file for the networking/load balancing
emulation. This means when you can create things like DNS records the same way you would with the real cloud services that your `/etc/hosts` file will be updated accordingly.

TODO: Provide instructions for downloading, installing and running the host agent.

**Configuration**

The host agent shares exactly the same configuration as the main server, see the [configuration](#configuration) section above.

### Running Directly On The Host

TODO: Provide instructions for downloading and running the binary locally.

## Google Cloud Service Endpoints

Cloud::1 provides some google cloud services that are accessible via a HTTP API along with a subset of services
that support gRPC.

All http endpoints are not secure unless you create your own reverse proxy that terminates with TLS, so use "http://" as the protocol for every endpoint.

*The square brackets represent the service name that can be used in configuration when selecting services to run.*

- [Secret Manager](https://cloud.google.com/secret-manager/docs/apis) [secretmanager] (HTTP, gRPC) - secretmanager.googleapis.local:5988/v1/
- [API Gateway](https://cloud.google.com/api-gateway/docs/apis) [apigateway] (HTTP) - apigateway.googleapis.local:5988/v1beta/
- [Cloud Storage](https://cloud.google.com/storage/docs/json_api) [storage] (HTTP) - storage.googleapis.local:5988/storage/v1/