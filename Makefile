BINARY_NAME=main.out

build:
	go build -o out/${BINARY_NAME} main.go

run:
	go build -o out/${BINARY_NAME} main.go
	./out/${BINARY_NAME}

test:
	go test -v ./tests/...

cover:
	go test -coverprofile cover.out
	go tool cover -html=cover.out

clean:
	go clean
	rm -rf ./out
	rm -rf ./cache
	rm *.out
