functions := $(shell cd src && find . -name \*main.go)

build: clean
	@mkdir -p ../out
	@for f in $(functions) ; do \
			dirname=$$(dirname $$f); \
			file=$$(basename $$f .go); \
      cd src && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../out/$$dirname/$$file $$f; \
    done

clean:
	@cd src && go clean
	@rm -rf ./out

zip:
	@for f in out/lambdas/*; do \
        zip -j $$f.zip $$f; \
    done

format:
	@cd src && gofmt -s -w .
