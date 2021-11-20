# Debugging

This document provides helpful tips and pointers for debugging issues
when contributing to the cloud uno project.

## Docker network issues when running the server locally in Docker

```
ERROR: Pool overlaps with other one on this address space
```
If you see this error when running the server and host agent locally,
run `docker network prune` to clean out dangling networks that use the same
Docker subnet.
