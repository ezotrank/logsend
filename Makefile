.PHONY: all nuke

all:
	go get github.com/mattn/gom
	mkdir -p vendor/bin
	gom install
	gofmt -w ./logsend ./main.go
	gom test ./logsend
	go vet ./logsend
	gom build -o vendor/bin/logsend ./main.go
	GOOS=linux gom build -o vendor/bin/logsend_linux ./main.go

nuke:
	go clean -i