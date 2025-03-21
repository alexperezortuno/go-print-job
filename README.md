# go-print-job

## Run

```shell
go run main.go
```

## Build

```shell
CGO_ENABLED=1 GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) go build -ldflags "-X 'main.username=your_username' -X 'main.password=your_password' -X 'main.urlLogin=http://example.com/login' -X 'main.urlJobs=http://example.com/job'" -o $(pwd)/dist/PrintJobsMonitor_$(go env GOOS)_$(go env GOARCH).exe
```

```shell
cd dist
PrintJobsMonitor_$(go env GOOS)_$(go env GOARCH).exe
```
