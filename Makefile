.PHONY: all nuke build deploy deps test

deploy_user := $(DEPLOY_USER)
deploy_target := $(DEPLOY_TO)

all: deps test build

deps:
	go get github.com/mattn/gom
	mkdir -p vendor/bin
	gom install

test:
	gofmt -w ./logsend ./main.go
	gom test ./logsend
	go vet ./logsend

build:
	gom build -o vendor/bin/logsend ./main.go
	GOOS=linux gom build -o vendor/bin/logsend_linux ./main.go

deploy: build
	ssh ${deploy_user}@${deploy_target} "mkdir -p ~/logsend && touch ~/logsend/config.json"
	scp vendor/bin/logsend_linux ${deploy_user}@${deploy_target}:"~/logsend"

nuke:
	go clean -i

