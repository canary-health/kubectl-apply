.DEFAULT_GOAL := all

########################
## Internal variables
########################
ORG_NAME=canary-health
APP_NAME=kubectl-apply
AWS_REGION=us-east-1

########################
## External Variables
########################
$(shell echo APP_NAME=${APP_NAME} > .env)
$(shell echo AWS_REGION=${AWS_REGION} >> .env)

########################
## Helpers variables
########################
M=$(shell printf "\033[34;1m▶\033[0m")
TIMESTAMP := $(shell /bin/date "+%Y-%m-%d_%H-%M-%S")

######
## Build targets
######

build-cli: dep ; $(info $(M) Building project cli...)
	CGO_ENABLED=0 GOOS=linux retool do go build -o ./bin/$(APP_NAME) ./

######
## Setup commands
######
.PHONY: setup dep build-server run-server

setup: ; $(info $(M) Fetching github.com/twitchtv/retool...)
	go get -u github.com/twitchtv/retool
	
dep: setup ; $(info $(M) Ensuring vendored dependencies are up-to-date...)
	retool sync

######
## Docker commands
######
.PHONY: build-image run-container

build-image: ; $(info $(M) Building docker image...)
	docker build --file ./build/docker/Dockerfile --tag $(APP_NAME) .

run-container: ; $(info $(M) Running docker container...)
	docker run -p $(PORTS) $(APP_NAME)


######
## Test commands
######
.PHONY: test test-coverage test-coverage-html

test: rm-coverage; $(info $(M) Running application tests...)
	go test ./... -cover -covermode=count -coverprofile=coverage.out

test-coverage: ;
	go tool cover -func=coverage.out

test-coverage-html: ;
	go tool cover -html=coverage.out

######
## Clean up commands
######
.PHONY: clean rm-bin rm-pb rm-tools rm-coverage

clean: rm-bin rm-pb rm-tools rm-coverage; $(info $(M) Removing ALL generated files... )

rm-bin: ; $(info $(M) Removing ./bin files... )
	sudo rm -rf ./bin

rm-tools: ; $(info $(M) Removing ./_tools files... )
	sudo rm -rf ./_tools

rm-coverage: ; $(info $(M) Removing coverage.out files... )
	$(RM) ./coverage.out
