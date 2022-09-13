BINARY_NAME=main.out

build:
	go build -o out/${BINARY_NAME} main.go

run:
	go build -o out/${BINARY_NAME} main.go
	./out/${BINARY_NAME}

clean:
	go clean
	rm -rf ./out
	rm -rf ./cache
