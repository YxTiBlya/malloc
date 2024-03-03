run-go:
	@go build -o ./bin/go.exe ./go/main.go && cd bin && go.exe

run-cgo:
	@go build -o ./bin/cgo.exe ./cgo && cd bin && cgo.exe
