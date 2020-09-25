docker_itest: Dockerfile.itest
	docker build -t go-jsonpb/itest -f Dockerfile.itest .
	docker run --rm -v $(shell pwd):/go/protoc-gen-go-json go-jsonpb/itest make itest

test:
	go test -v ./...

itest:
	go test -tags itest -v ./...
