all:
	go build -o cli cmd/main.go && go build -o server cmd/server.go

server:
	./server

test:
	go test -v -count=1 -race ./...

clean:
	rm -rf cli server
