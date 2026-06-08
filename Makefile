all:
	go build -o cli cmd/main.go && go build -o server cmd/server.go

server:
	./server

clean:
	rm -rf cli server
