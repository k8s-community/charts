all: push

BUILDTAGS=

# Use the 0.0.0 tag for testing, it shouldn't clobber any release builds
APP=charts
PROJECT?=charts
RELEASE?=0.0.9
GOOS?=linux
REPOSITORY?=community-charts
REGISTRY?=containers.k8s.community
CHARTS_SERVICE_PORT?=8080
CHARTS_SERVICE_HEALTH_PORT?=8082

NAMESPACE?=dev
PREFIX?=${REGISTRY}/${NAMESPACE}/${APP}
CONTAINER_NAME?=${APP}-${NAMESPACE}

ifeq ($(NAMESPACE), default)
  PREFIX=${REGISTRY}/${APP}
  CONTAINER_NAME=${APP}
endif

REPO_INFO=$(shell git config --get remote.origin.url)

ifndef COMMIT
  COMMIT := git-$(shell git rev-parse --short HEAD)
endif

vendor: clean
	go get -u github.com/Masterminds/glide \
	&& glide install

utils:
	go get -u github.com/golang/lint/golint

build: vendor
	CGO_ENABLED=0 GOOS=${GOOS} go build -a -installsuffix cgo \
		-ldflags "-s -w -X ${PROJECT}/version.RELEASE=${RELEASE} -X ${PROJECT}/version.COMMIT=${COMMIT} -X ${PROJECT}/version.REPO=${REPO_INFO}" \
		-o ${APP}-server

container: build
	docker build --pull -t $(PREFIX):$(RELEASE) .

push: container
	docker push $(PREFIX):$(RELEASE)

run: container
	docker run --name ${CONTAINER_NAME} -p ${CHARTS_SERVICE_PORT}:${CHARTS_SERVICE_PORT} \
	  -p ${CHARTS_SERVICE_HEALTH_PORT}:${CHARTS_SERVICE_HEALTH_PORT} \
	  -e "CHARTS_SERVICE_PORT=${CHARTS_SERVICE_PORT}" \
	  -e "CHARTS_SERVICE_HEALTH_PORT=${CHARTS_SERVICE_HEALTH_PORT}" \
	  -d $(PREFIX):$(RELEASE)

deploy: push
	helm repo up \
    && helm upgrade ${CONTAINER_NAME} ${REPOSITORY}/${APP} --namespace ${NAMESPACE} --set image.tag=${RELEASE} -i --wait

fmt:
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"gofmt -s -l {{.Dir}}"{{end}}' $(shell go list ${PROJECT}/... | grep -v vendor) | xargs -L 1 sh -c

lint: utils
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"golint {{.Dir}}/..."{{end}}' $(shell go list ${PROJECT}/... | grep -v vendor) | xargs -L 1 sh -c

vet:
	@echo "+ $@"
	@go vet $(shell go list ${PROJECT}/... | grep -v vendor)

test: vendor utils fmt lint vet
	@echo "+ $@"
	@go test -v -race -tags "$(BUILDTAGS) cgo" $(shell go list ${PROJECT}/... | grep -v vendor)

cover:
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}}"{{end}}' $(shell go list ${PROJECT}/... | grep -v vendor) | xargs -L 1 sh -c

clean:
	rm -f ${APP}-server
