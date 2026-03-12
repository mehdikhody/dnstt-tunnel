PLATFORMS := linux/amd64 linux/arm64 windows/amd64

build:
	go build -o build/dnstt-tunnel-windows-amd64.exe

run:
	make build
	build\dnstt-tunnel-windows-amd64.exe