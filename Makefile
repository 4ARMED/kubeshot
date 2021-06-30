EXECUTABLE := "kubeshot"
GITVERSION := $(shell git describe --dirty --always --tags --long)
PACKAGENAME := $(shell go list -m -f '{{.Path}}')

build: clean test
	go build -ldflags "-extldflags '-static' -X ${PACKAGENAME}/internal/config.GitVersion=${GITVERSION}" -o ${EXECUTABLE} .

build-quick: clean
	go build -ldflags "-extldflags '-static' -X ${PACKAGENAME}/internal/config.GitVersion=${GITVERSION}" -o ${EXECUTABLE} .

build-linux:
	GOOS=linux go build -ldflags "-extldflags '-static' -X ${PACKAGENAME}/internal/config.GitVersion=${GITVERSION}" -o ${EXECUTABLE}-linux

install: clean
	go install -a -tags netgo -ldflags "-w -extldflags '-static' -X ${PACKAGENAME}/internal/config.GitVersion=${GITVERSION}"

clean:
	@rm -f ${EXECUTABLE}

test:
	go test -v ./...

docker:
	docker build -f build/Dockerfile . -t ${EXECUTABLE}