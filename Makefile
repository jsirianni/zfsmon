ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

UNIXTIME := $(shell date +'%s')
VERSION := $(shell cat main.go | grep "const VERSION" | cut -c 17- | tr -d '"')

$(shell mkdir -p artifacts)

build: clean
	$(info building zfsmon ${VERSION})

	@docker build \
	    --build-arg version=${VERSION} \
	    -t zfsmon:${VERSION} .

	@docker create -ti --name ${UNIXTIME} zfsmon:${VERSION} bash && \
		docker cp ${UNIXTIME}:/zfsmon/artifacts/. artifacts/

	# cleanup
	@docker rm -fv ${UNIXTIME} &> /dev/null

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
