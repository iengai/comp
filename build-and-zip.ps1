# build-and-zip.ps1
$env:GOARCH="arm64"
$env:GOOS="linux"
go build -ldflags="-s -w" -o bin/bootstrap functions/hello-world/main.go
Compress-Archive -Path bin/bootstrap -DestinationPath bin/hello-world.zip -Update