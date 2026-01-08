.PHONY: test lint fmt tidy tag

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

lint:
	go vet ./...

tag:
	@test -n "$(v)" || (echo "Usage: make tag v=v0.1.0" && exit 1)
	git tag $(v)
	git push origin $(v)
