.SILENT:


BIN := main.exe

build:
	go build -o ${BIN} cmd/main.go

test:
	go test ./...

clean:
	go clean
	rm -f ${BIN}
