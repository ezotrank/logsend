# Examples:
# DEPLOY_USER=user DEPLOY_CONFIG=some_config DEPLOY_TO="host1 host2" make deploy

.PHONY: all nuke build deploy deps test

deploy_user := $(DEPLOY_USER)
deploy_target := $(DEPLOY_TO)
deploy_config := $(DEPLOY_CONFIG)

all: deps format test build

deps:
	go get

test:
	go test ./logsend
	go vet ./logsend

format:
	gofmt -w ./logsend ./main.go

build:
	go build -o $$GOPATH/bin/logsend ./main.go
	GOOS=linux go build -o $$GOPATH/bin/logsend_linux ./main.go

deploy: build deploy_copy deploy_update_config deploy_restart

deploy_copy:
	for host in ${deploy_target}; do \
		ssh ${deploy_user}@$$host "mkdir -p ~/logsend" ; \
		scp $$GOPATH/bin/logsend_linux ${deploy_user}@$$host:"~/logsend/logsend_linux.NEW" ; \
		ssh ${deploy_user}@$$host "cd ~/logsend && mv logsend_linux.NEW logsend" ; \
	done


deploy_update_config:
	for host in ${deploy_target}; do \
		scp ./own_configs/${deploy_config}_config.json ${deploy_user}@$$host:"~/logsend/config.json" ; \
	done


deploy_restart:
	for host in ${deploy_target}; do \
		ssh ${deploy_user}@$$host "cd ~/logsend && cat logsend.pid |xargs kill || true" ; \
	done


nuke:
	go clean -i

