.PHONY: build

build:
	CGO=0 GOOS=js GOARCH=wasm go build -o ../../../public/wasm/activity-service.wasm -ldflags="-s -w" -trimpath