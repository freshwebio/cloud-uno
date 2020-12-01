# Cloud::1

Cloud::1 enables running cloud services on your local machine for faster, cheaper, end-to-end development.

*If you were wondering, the name Cloud::1 (cloud one) is derived from the ::1 IPv6 address that is the shorthand for the loopback localhost ip.*

## Prerequisites

### Running as a container or standalone

- [Docker](https://docs.docker.com/get-docker/) - Docker is required to run Cloud::1 as a container. Some services in Cloud::1 run in their own docker containers so even when you run the app as a standalone binary, you still need Docker.

### Build yourself

- [Docker](https://docs.docker.com/get-docker/) - Some services in Cloud::1 run in their own docker containers so even when you run the app as a standalone binary, you still need Docker.
- [Go](https://golang.org/dl/) - Cloud::1 is implemented in Go so in order to build from source you need Go installed.

## Setting up

### Local service resolution

If you want to use Cloud::1 for Google Cloud or Azure you need to configure your machine so Cloud::1 knows which local service to forward to.

**You can skip this section if all you want to use Cloud::1 for is AWS services.**

You can use the helper script to sort this out for you:
TODO: Create and document helper script

To do this manually, open up your hosts file `/etc/hosts` for linux or macos machines and `C:\Windows\System32\Drivers\etc\hosts` on windows.

Then add this to the end of the file:
```
googleapis.local 127.0.0.1
microsoft.local 127.0.0.1
```

Now save the file.

### Running in Docker

TODO: Provide instructions for running as a Docker container and with Docker compose.

### Running standalone

TODO: Provide instructions for downloading and running the binary locally.

## Google Cloud Service endpoints

- [Secret Manager](https://cloud.google.com/secret-manager/docs/apis) - secretmanager.googleapis.local:5988