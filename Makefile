SOURCE = bank.go
CONTAINER = plaid_bank_sym
SOURCE_PATH = /go/src/github.com/buger/$(CONTAINER)
TEST = .
DRUN = docker run -v `pwd`:$(SOURCE_PATH) -i -t $(CONTAINER)

build:
	docker build -t $(CONTAINER) .

race:
	$(DRUN) --env GORACE="halt_on_error=1" go test ./. $(ARGS) -v -race -timeout 15s

test:
	$(DRUN) go test $(LDFLAGS) ./ -run $(TEST) -timeout 10s $(ARGS) -v

fmt:
	$(DRUN) go fmt ./...

vet:
	$(DRUN) go vet ./.

run:
	$(DRUN) go run $(SOURCE)

bash:
	$(DRUN) /bin/bash