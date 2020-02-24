ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

UNIXTIME := $(shell date +'%s')
VERSION := $(shell cat main.go | grep "const VERSION" | cut -c 17- | tr -d '"')
WORK_DIR := $(shell pwd)

$(shell mkdir -p artifacts)

build: build-zfs-7

build-zfs-7:
	$(info building zfsmon ${VERSION})

	@cp build/zfs7/go.mod.zfs.7 go.mod
	@cp build/zfs7/go.sum.zfs.7 go.sum
	@cp build/zfs7/Dockerfile Dockerfile

	@docker build \
	    --build-arg version=${VERSION}-zfs-7 \
	    -t zfsmon:${VERSION} .

	@docker create -ti --name ${UNIXTIME} zfsmon:${VERSION} bash && \
		docker cp ${UNIXTIME}:/zfsmon/artifacts/. artifacts/ && \
		docker cp ${UNIXTIME}:/zfsmon/go.mod build/zfs7/go.mod.zfs.7 && \
		docker cp ${UNIXTIME}:/zfsmon/go.sum build/zfs7/go.sum.zfs.7 && \
		docker rm -fv ${UNIXTIME}

	@rm -f go.mod go.sum Dockerfile

test:
	go test ./...

lint:
	golint ./...

fmt:
	go fmt ./...

clean:
	@rm -f go.mod go.sum Dockerfile
	@rm -rf artifacts/*

quick:
	go build -o zfsmon

prune-docker:
	docker system prune --force
