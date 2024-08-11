# Large Scale Workshop

This project displays a basic distributed, service-oriented system including the following concepts: Remote Procedure Call, Service Discovery, Distributed cache and Message Queue.

The system has three sevices:
- Test Service - provides several functionalities
- Registry Service - provides service discovery
- Cache Service - provides in-memory cache using Chord Distributed Hash Table


## Getting started
- `git clone`
- load within the project dev container
- the image mounts the large-scale-workshop directory into: `/workspaces/<cloned-repo-name>/`
- our "module name" is `github.com/TAULargeScaleWorkshop/RLAD`

## Building
In the root directory run `./build.sh` that installs required dependencies and builds the app to `./output/`.

## Running
Run the app using `./output/start.sh` that starts 3 services of each type: Registry, Cache and Test.

## Testing
TestService:
```
cd /workspaces/RLAD/services/test-service/client/
go test -v
```
Chord DHT:
```
cd /workspaces/RLAD/services/reg-service/servant/dht
go test -v
```