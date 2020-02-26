GOCMD=go
GOBUILD=$(GOCMD) build -o ./.aws-sam/build/
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
SAMBUILD=./.aws-sam/build/

.PHONY: deps clean build

all: deps test clean build

# test all packages
test:
	$(GOTEST) ./...

# get all packages
deps:
	$(GOGET) -u ./...

# static analysis
static:
	$(GOVET) ./...

clean: 
	rm -rf $(SAMBUILD)PrismaAlertDispatcher/dispatcher
	rm -rf $(SAMBUILD)PrismaAlertNotification/snsalert
	rm -rf $(SAMBUILD)PrismaVPCKiller/vpckiller
	rm -rf $(SAMBUILD)PrismaFalseAlertRemover/remover
	rm -rf $(SAMBUILD)PrismaOnboarding/onboarding
	
build:
	GOOS=linux GOARCH=amd64 $(GOBUILD)PrismaAlertDispatcher/dispatcher ./remediation/dispatcher/main.go
	GOOS=linux GOARCH=amd64 $(GOBUILD)PrismaAlertNotification/snsalert ./remediation/snsalert/snsalert.go
	GOOS=linux GOARCH=amd64 $(GOBUILD)PrismaVPCKiller/vpckiller ./remediation/vpckiller/vpckiller.go
	GOOS=linux GOARCH=amd64 $(GOBUILD)PrismaFalseAlertRemover/remover ./remediation/falsealert/remover.go
	GOOS=linux GOARCH=amd64 $(GOBUILD)PrismaOnboarding/onboarding ./remediation/onboarding/onboarding.go