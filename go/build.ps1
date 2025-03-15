$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o aipipe.exe ./cmd/aipipe
Write-Output "Build completed: windows/amd64"
