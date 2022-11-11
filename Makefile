BINARY=printer
BUILD=$$(vtag --no-meta)
TAG="subtlepseudonym/${BINARY}:${BUILD}"

default: build

build: format
	go build -o ${BINARY} -v ./main.go

docker: format
	docker build --network=host --tag ${TAG} -f Dockerfile .

test: format
	gotest --race ./...

format fmt:
	gofmt -l -w -e .

clean:
	go mod tidy
	go clean
	rm -f $(BINARY)

get-tag:
	echo ${BUILD}

.PHONY: all build format fmt clean get-tag
