# RLAD

Large Scale Workshop

## General Notes
- the image mounts the large-scale-workshop directory into: `/workspaces/<cloned-repo-name>/`
- some necessary dependencies for python 3.11 are not included in the base docker image, see extra installations.
- our "module name" is `github.com/TAULargeScaleWorkshop/RLAD`

## Section 1 - running main.go
```
go get
go build -o ./output/large-scale-workshop
./output/large-scale-workshop ./services/test-service/service/TestService.yaml
```

## Section 2 - running the test
```
cd /workspaces/large-scale-workshop/interop/
go clean -testcache
go test -v -tags=interop
```

## Section 3
First, compiling IDL into Go protocol buffers code and gRPC to generate `TestService.pb.go` and `TestService_grpc.pb.go`
```
cd /workspaces/large-scale-workshop/services/test-service/common
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
cd /workspaces/large-scale-workshop/services/test-service/client/
go test -v
```

#### Extra installations
```
sudo apt-get update && sudo apt-get install -y python3.11-dev
python3.11 -m pip install beautifulsoup4 requests
```

## Section 4
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

make sure you get the correct version of MetaFFI, run at the root directory of the project and build everything:
```
go get github.com/MetaFFI/lang-plugin-go@v0.1.2
go mod tidy
go build -o ./output/large-scale-workshop
```
We have three components now that should run in separate terminals:
1. RegService:
```
./output/large-scale-workshop ./services/reg-service/service/RegService.yaml
```
2. TestService:
```
./output/large-scale-workshop ./services/test-service/service/TestService.yaml
```
3. TestServiceClient:
```
cd /workspaces/large-scale-workshop/services/test-service/client/
go test -v
```
