# RLAD

Large Scale Workshop

## General Notes
- the image mounts the large-scale-workshop directory into: `/workspaces/large-scale-workshop/`
- some necessary dependencies for python 3.11 are not included in the base docker image, see extra installations.
- our "module name" is `github.com/TAULargeScaleWorkshop/RLAD`

### Section 1 - running main.go
```
go get
go build -o ./output/large-scale-workshop
./output/large-scale-workshop ./services/test-service/service/TestService.yaml
```

### Section 2 - running the test
```
cd /workspaces/large-scale-workshop/interop/
go clean -testcache
go test -v -tags=interop
```

### Section 3
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
Next:
```
go get
go build -o ./output/large-scale-workshop
./output/large-scale-workshop ./services/test-service/service/TestService.yaml
cd services/test-service/client/
go test -v
```

#### Extra installations
```
sudo apt-get update && sudo apt-get install -y python3.11-dev
python3.11 -m pip install beautifulsoup4 requests
```
