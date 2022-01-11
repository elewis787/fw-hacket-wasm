build:
	GOOS=js GOARCH=wasm go build -o ./bin/hacket.wasm *.go
