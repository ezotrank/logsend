.PHONY: all nuke build deploy deps test

deploy_user := $(DEPLOY_USER)
deploy_target := $(DEPLOY_TO)
deploy_config := $(DEPLOY_CONFIG)

all: deps test build

deps:
	go get

test:
	gofmt -w ./logsend ./main.go
	go test ./logsend
	go vet ./logsend

build:
	go build -o $$GOPATH/bin/logsend ./main.go
	GOOS=linux go build -o $$GOPATH/bin/logsend_linux ./main.go

deploy: build deploy_copy deploy_update_config deploy_restart

deploy_copy:
	for host in ${deploy_target}; do \
		ssh ${deploy_user}@$$host "mkdir -p ~/logsend" ; \
		scp vendor/bin/logsend_linux ${deploy_user}@$$host:"~/logsend/logsend_linux.NEW" ; \
		ssh ${deploy_user}@$$host "cd ~/logsend && mv logsend_linux.NEW logsend_linux" ; \
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

