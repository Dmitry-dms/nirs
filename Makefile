.SILENT:


run:
	go run cmd/main.go

build:
	go build -o . cmd/main.go
t:
	go run test/main.go