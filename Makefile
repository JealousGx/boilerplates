functions := $(shell cd pkg && find . -name \*main.go)

build: clean
	@mkdir -p ../out
	@for f in $(functions) ; do \
			dirname=$$(dirname $$f); \
			file=$$(basename $$f .go); \
      cd pkg && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../out/$$dirname/$$file $$f; \
    done

clean:
	@cd pkg && go clean
	@rm -rf ./out

zip:
	@for f in out/lambdas/*; do \
        zip -j $$f.zip $$f; \
    done

format:
	@cd pkg && gofmt -s -w .
