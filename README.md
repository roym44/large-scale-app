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

## Section 3
First, compiling IDL into Go protocol buffers code and gRPC to generate `TestService.pb.go` and `TestService_grpc.pb.go`
```
cd /workspaces/RLAD/services/test-service/common
protoc -I=. \
       --go_out=. \
       --go_opt=paths=source_relative \
       --go_opt=MTestService.proto=github.com/TAULargeScaleWorkshop/RLAD/services/test-service \
       --go-grpc_out=. \
       --go-grpc_opt=paths=source_relative \
       --go-grpc_opt=MTestService.proto=github.com/TAULargeScaleWorkshop/RLAD/services/testservice/TestService.proto \
       TestService.proto
```
Next, to run the server:
```
go get
go build -o ./output/large-scale-workshop
./output/large-scale-workshop ./services/test-service/service/TestService.yaml
```
And test the client:
```
cd /workspaces/RLAD/services/test-service/client/
go test -v
```

## Section 4
### Cluster & Registry
First, compiling IDL into Go protocol buffers code and gRPC to generate `RegService.pb.go` and `RegService_grpc.pb.go`
```
cd /workspaces/RLAD/services/reg-service/common
protoc -I=. \
       --go_out=. \
       --go_opt=paths=source_relative \
       --go_opt=MRegService.proto=github.com/TAULargeScaleWorkshop/RLAD/services/reg-service \
       --go-grpc_out=. \
       --go-grpc_opt=paths=source_relative \
       --go-grpc_opt=MRegService.proto=github.com/TAULargeScaleWorkshop/RLAD/services/regservice/RegService.proto \
       RegService.proto
```
Now compiling our main:
```
cd /workspaces/RLAD
go get
go build -o ./output/large-scale-workshop
```
We have three components now that should run in separate terminals:
1. RegService:
```
./output/large-scale-workshop ./services/reg-service/service/RegService.yaml
```
Unit testing for RegService:
```
cd /workspaces/RLAD/services/reg-service/client/
go test -v
```
2. TestService:
```
./output/large-scale-workshop ./services/test-service/service/TestService.yaml
```
3. TestServiceClient:
```
cd /workspaces/RLAD/services/test-service/client/
go test -v
```


### Cluster Registry Service & Cache Service
Make sure you get the correct version of MetaFFI, run at the root directory of the project and build everything:
```
cd /workspaces/RLAD
go get github.com/MetaFFI/lang-plugin-go@v0.1.2
go mod tidy
go build -o ./output/large-scale-workshop
```
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
1. RegService:
```
cd /workspaces/RLAD/services/reg-service/servant
/workspaces/RLAD/output/large-scale-workshop /workspaces/RLAD/services/reg-service/service/RegService.yaml
```
2. TestService:
```
cd /workspaces/RLAD
/workspaces/RLAD/output/large-scale-workshop /workspaces/RLAD/services/test-service/service/TestService.yaml
```
3. TestServiceClient:
```
cd /workspaces/RLAD/services/test-service/client/
go test -v
```
