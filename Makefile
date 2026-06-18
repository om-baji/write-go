all:
	go build -o cli cmd/cli/main.go && go build -o server cmd/server/main.go

server:
	./server

test:
	go test -v -count=1 -race ./...

clean:
	rm -rf cli server
