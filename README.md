# RLAD

Large Scale Workshop

## General Notes
- the image mounts the large-scale-workshop directory into: `/workspaces/large-scale-workshop/`

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

#### Extra installations
```
sudo apt-get update && sudo apt-get install -y python3.11-dev
python3.11 -m pip install beautifulsoup4 requests
```
