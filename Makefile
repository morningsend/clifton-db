GOBIN=$(CURDIR)/bin
OUT=cliftondbctl

all:
	go build -o ${GOBIN}/${OUT} cliftondb/cmd/*.go

generate:
	go generate $(CURDIR)/pkg/storagepb/*.go