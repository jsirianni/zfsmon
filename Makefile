ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

UNIXTIME := $(shell date +'%s')
VERSION := $(shell cat main.go | grep "const VERSION" | cut -c 17- | tr -d '"')
WORK_DIR := $(shell pwd)

$(shell mkdir -p artifacts)

build: build-zfs-7

build-zfs-7: clean
	$(info building zfsmon ${VERSION})

	@cp go.mod.zfs.7 go.mod

	@docker build \
	    --build-arg version=${VERSION}-zfs-7 \
	    -t zfsmon:${VERSION} \
		-f Dockerfile.zfs.7 .

	@docker create -ti --name ${UNIXTIME} zfsmon:${VERSION} bash && \
		docker cp ${UNIXTIME}:/zfsmon/artifacts/. artifacts/ && \
		docker cp ${UNIXTIME}:/zfsmon/go.mod go.mod.zfs.7 && \
		docker cp ${UNIXTIME}:/zfsmon/go.sum go.sum.zfs.7 && \
		docker rm -fv ${UNIXTIME}

test:
	go test ./...

lint:
	golint ./...

fmt:
	go fmt ./...

clean:
	$(shell rm -f go.mod go.sum)
	$(shell rm -rf artifacts/*)

quick:
	go build -o zfsmon

prune-docker:
	docker system prune --force
