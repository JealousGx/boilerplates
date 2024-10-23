functions := $(shell cd pkg/lambdas && find . -name \*main.go)

build: clean
	@mkdir -p ../out
	@for f in $(functions) ; do \
			dirname=$$(dirname $$f); \
			dirname=$$(basename $$dirname); \
			file=$$(basename $$f .go); \
			echo "Building $$f in directory $$dirname"; \
			(cd pkg/lambdas && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../out/lambdas/$$dirname $$f); \
    done

clean:
	@cd pkg && go clean
	@rm -rf ./out

zip:
	@for f in out/lambdas/*; do \
			if [ -f $$f ]; then \
					zip -j $$f.zip $$f; \
			else \
					echo "Skipping $$f as it does not exist or is empty"; \
			fi \
	done

zip-test:
	@for f in out/lambdas/*; do \
			zip -j $$f.zip $$f; \
  	done

format:
	@cd pkg && gofmt -s -w .
