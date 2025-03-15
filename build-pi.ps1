$env:GOOS = "linux"
$env:GOARCH = "arm64"
go build -o aipipe-linux-arm64 ./cmd/aipipe
Write-Host "Build completed: linux/arm64"