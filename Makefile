ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

UNIXTIME := $(shell date +'%s')
VERSION := $(shell cat main.go | grep "const VERSION" | cut -c 17- | tr -d '"')
WORK_DIR := $(shell pwd)

$(shell mkdir -p artifacts)

build: clean
	$(info building zfsmon ${VERSION})

	@docker build \
	    --build-arg version=${VERSION} \
	    -t zfsmon:${VERSION} .

	@docker create -ti --name ${UNIXTIME} zfsmon:${VERSION} bash && \
		docker cp ${UNIXTIME}:/zfsmon/artifacts/. artifacts/ && \
		docker cp ${UNIXTIME}:/zfsmon/go.mod go.mod && \
		docker cp ${UNIXTIME}:/zfsmon/go.sum go.sum && \
		docker rm -fv ${UNIXTIME}

test:
	go test ./...

lint:
	golint ./...

fmt:
	go fmt ./...

clean:
	$(shell rm -rf artifacts/*)

quick:
	go build -o zfsmon

prune-docker:
	docker system prune --force
