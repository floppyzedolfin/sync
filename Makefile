PROJECT=sync
VERSION=1.0.0

# GO111MODULE=on means we're using modules and not the olde vendor dir and $GO*** environment variables
# CGO_ENABLED=0 means we're not looking for C libs when using the network packages (which comes in handy when using scratch images)
GO=GO111MODULE=on CGO_ENABLED=0 go
PLAYGROUND=${PWD}/example/playground
REPLICA=${PWD}/example/replica

SERVER_ADDRESS=178.18.0.23
SERVER_PORT=8405
LOCAL_REPLICA_DIR=/data/replica

generate:
	${GO} generate ./...

compile-watcher:
	${GO} build -o build/watcher/watcher.out ./watcher/cmd

compile-replica:
	${GO} build -o build/replica/replica.out ./replica/cmd

compile-all: compile-watcher compile-replica

test:
	${GO} test ./...

build-all: compile-all
	docker build -t ${PROJECT}-replica:${VERSION} build/replica

network:
	docker network create --subnet=178.18.0.0/16 network-${PROJECT} || true

run-replica: network
	docker run -d -p ${SERVER_PORT}:${SERVER_PORT} -v ${REPLICA}:${LOCAL_REPLICA_DIR} --network network-${PROJECT} --ip ${SERVER_ADDRESS} --rm --name replica ${PROJECT}-replica:${VERSION} -- -local_replica ${LOCAL_REPLICA_DIR} -port ${SERVER_PORT}

make run-watcher:
	build/watcher/watcher.out --watched_dir ${PLAYGROUND} --server_addr ${SERVER_ADDRESS}:${SERVER_PORT}

run-all: run-watcher run-replica

stop-all:
	docker stop replica
