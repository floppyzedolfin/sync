PROJECT=sync
VERSION=1.0.0

# GO111MODULE=on means we're using modules and not the olde vendor dir and $GO*** environment variables
# CGO_ENABLED=0 means we're not looking for C libs when using the network packages (which comes in handy when using scratch images)
GO=GO111MODULE=on CGO_ENABLED=0 go
PLAYGROUND=${PWD}/example/playground
REPLICA=${PWD}/example/replica

generate:
	${GO} generate ./...

compile-reference:
	${GO} build -o build/reference/reference.out ./reference/service/cmd

compile-watcher:
	${GO} build -o build/watcher/watcher.out ./watcher/cmd

compile-all: compile-reference compile-watcher

test:
	${GO} test ./...

build-all: compile
	docker build -t ${PROJECT}-reference:${VERSION} build/reference

network:
	docker network create --subnet=178.18.0.0/16 network-${PROJECT} || true

run-reference: network
	docker run -d -p 8405:8405 -v ${REPLICA}:/data/replica --network network-${PROJECT} --ip 178.18.0.23 --rm --name reference ${PROJECT}-reference:${VERSION}

make run-watcher:
	build/watcher/watcher.out ${PLAYGROUND}

run-all: run-reference run-watcher

stop-all:
	docker stop reference
