run-go:
	@go build -o ./bin/go.exe ./go/main.go && cd bin && go.exe

run-cgo:
	@go build -o ./bin/cgo.exe ./cgo && cd bin && cgo.exe

run-c:
	@gcc ./c/y_malloc.c -o ./bin/c.exe && cd bin && c.exe
