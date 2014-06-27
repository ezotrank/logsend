.PHONY: all nuke build deploy deps test

deploy_user := $(DEPLOY_USER)
deploy_target := $(DEPLOY_TO)
deploy_config := $(DEPLOY_CONFIG)

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

deploy: build deploy_copy deploy_update_config deploy_restart

deploy_copy:
	ssh ${deploy_user}@${deploy_target} "mkdir -p ~/logsend"
	scp vendor/bin/logsend_linux ${deploy_user}@${deploy_target}:"~/logsend/logsend_linux.NEW"
	ssh ${deploy_user}@${deploy_target} "cd ~/logsend && mv logsend_linux.NEW logsend_linux"

deploy_update_config:
	scp ./own_configs/${deploy_config}_config.json ${deploy_user}@${deploy_target}:"~/logsend/config.json"

deploy_restart:
	ssh ${deploy_user}@${deploy_target} "cd ~/logsend && cat logsend.pid |xargs kill || true"

nuke:
	go clean -i

