# RLAD

Large Scale Workshop

## General Notes
- the image mounts the large-scale-workshop directory into: `/workspaces/<cloned-repo-name>/`
- some necessary dependencies for python 3.11 are not included in the base docker image, see extra installations.
- our "module name" is `github.com/TAULargeScaleWorkshop/RLAD`

### Extra installations
```
sudo apt-get update && sudo apt-get install -y python3.11-dev
python3.11 -m pip install beautifulsoup4 requests
go get github.com/MetaFFI/lang-plugin-go@v0.1.2
go mod tidy
```

## Section 1 - running main.go
```
go get
go build -o ./output/large-scale-workshop
./output/large-scale-workshop ./services/test-service/service/TestService.yaml
```

## Section 2 - running the test
```
cd /workspaces/RLAD/interop/
go clean -testcache
go test -v -tags=interop
```

## Section 3 - TestService
First run: `/workspaces/RLAD/utils/testservice_proto.sh`\
Build server: `/workspaces/RLAD/build.sh`\
Run the server:
```
/workspaces/RLAD/output/large-scale-workshop /workspaces/RLAD/services/test-service/service/TestService.yaml
```
Test the client:
```
cd /workspaces/RLAD/services/test-service/client/
go test -v
```

## Section 4
### Cluster & Registry
First run: `/workspaces/RLAD/utils/regservice_proto.sh`\
Build: `/workspaces/RLAD/build.sh`\
We have three components now that should run in separate terminals:
1. RegService: `/workspaces/RLAD/utils/run_reg_service.sh`\
Unit testing for RegService:
```
cd /workspaces/RLAD/services/reg-service/client/
go test -v
```
2. TestService: `/workspaces/RLAD/utils/run_test_service.sh`\
3. TestServiceClient:
```
cd /workspaces/RLAD/services/test-service/client/
go test -v
```

### Cluster Registry Service & Cache Service
Chord DHT fixes:
- replace `Chord.class`
- `mv /workspaces/RLAD/files/xllr.openjdk.so /usr/local/metaffi/xllr.openjdk.so`
- `chmod 777 /usr/local/metaffi/xllr.openjdk.so`
We have the Chord DHT test:
```
cd /workspaces/RLAD/services/reg-service/servant/dht
go test -v
```
We have three components now that should run in separate terminals:
1. RegService: `/workspaces/RLAD/utils/run_reg_service.sh`
2. TestService: `/workspaces/RLAD/utils/run_test_service.sh`
3. CacheService: `/workspaces/RLAD/utils/run_cache_service.sh`
4. TestServiceClient:
```
cd /workspaces/RLAD/services/test-service/client/
go test -v
```
