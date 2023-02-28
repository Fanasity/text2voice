GO111MODULE=on

docker:
	go build -o text2voice main.go
	docker build -t text2voice:v0.1 .
