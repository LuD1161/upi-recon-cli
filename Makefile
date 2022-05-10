APP=upi-recon-cli
# https://medium.com/the-go-journey/adding-version-information-to-go-binaries-e1b79878f6f2
GIT_COMMIT=$(shell git rev-parse --short=10 HEAD)

.PHONY: build-and-execute
build-and-execute:
	go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" -o ${APP} main.go && chmod +x ./${APP} && ./${APP}

.PHONY: build
build:
	go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" -o ${APP} main.go

.PHONY: run
run:
	go run main.go

.PHONY: debug
debug: 
	export DEBUG=True && make build-and-execute
	
.PHONY: prod
prod: 
	export PROD=True && make build-and-execute
