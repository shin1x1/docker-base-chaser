.PHONY: test
test:
	go test -v ./...

.PHONY: vet
vet:
	go vet ./...


