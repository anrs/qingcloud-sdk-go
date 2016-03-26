
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

dep:
	go get -u github.com/nu7hatch/gouuidxs

test:
	go test -v ./...
