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
	  gzip -9 $$GOPATH/bin/logsend_linux ; \
		ssh ${deploy_user}@$$host "mkdir -p ~/logsend" ; \
		scp $$GOPATH/bin/logsend_linux.gz ${deploy_user}@$$host:"~/logsend/logsend.gz" ; \
		ssh ${deploy_user}@$$host "cd ~/logsend && gunzip -f logsend.gz" ; \
		rm -f $$GOPATH/bin/logsend_linux.gz ; \
	done


deploy_update_config:
	for host in ${deploy_target}; do \
		scp ./own_configs/${deploy_config}_config.json ${deploy_user}@$$host:"~/logsend/config.json" ; \
	done

deploy_update_monit:
	for host in ${deploy_target}; do \
		scp ./own_configs/${deploy_config}_monit.conf ${deploy_user}@$$host:"/etc/monit.d/logsend.conf" ; \
		ssh ${deploy_user}@$$host "sudo /etc/init.d/monit restart" ; \
	done


deploy_restart:
	for host in ${deploy_target}; do \
		ssh ${deploy_user}@$$host "cd ~/logsend && cat logsend.pid |xargs kill || true" ; \
	done


nuke:
	go clean -i

