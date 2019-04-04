cd ../../../
set GOPATH=%cd%
cd src/demo/idworker_srv
set GOARCH=amd64
set GOOS=windows
go build -v -ldflags="-s -w"