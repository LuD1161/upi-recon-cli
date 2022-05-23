APP=upi-recon-lambda
# https://medium.com/the-go-journey/adding-version-information-to-go-binaries-e1b79878f6f2
GIT_COMMIT=$(shell git rev-parse --short=10 HEAD)

.PHONY: build-and-execute
build-and-execute:
	chmod +x ./set-env-vars.sh && . ./set-env-vars.sh && go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" -o ${APP} *.go && chmod +x ./${APP} && ./${APP}

.PHONY: build
build:
	GOOS=linux CGO_ENABLED=0 go build -o main *.go
	GOOS=darwin CGO_ENABLED=0 go build -o ${APP} *.go

.PHONY: deploy
deploy:
	GOOS=linux CGO_ENABLED=0 go build -o main *.go
	zip -r function.zip main upi-recon-cli data/*
	mv function.zip ~/Desktop

.PHONY: debug
debug: 
	export DEBUG=True && make build-and-execute
	
.PHONY: prod
prod: 
	export PROD=True && make build-and-execute
