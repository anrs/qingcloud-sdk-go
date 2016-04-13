
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

dep:
	go get -u github.com/nu7hatch/gouuid

fmt:
	go fmt ./...

test: fmt
	go test -v ./...

cover:
	go test -cover -v ./...

coverprofile:
	go test -coverprofile=/tmp/c.out github.com/anrs/qingcloud-sdk-go/api
	go tool cover -html=/tmp/c.out
	go test -coverprofile=/tmp/c.out github.com/anrs/qingcloud-sdk-go/conn
	go tool cover -html=/tmp/c.out
