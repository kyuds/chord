# Makefile
# compile proto file and build go program.

GO_FILE = $(wildcard chord/*.go) $(wildcard *.go)
PROTO_FILE = pb/chord.proto
OUT_FILE = example

build: proto $(GO_FILE)
	go build -o $(OUT_FILE)

proto: $(PROTO_FILE)
	protoc --go_out=. --go_opt=paths=source_relative \
    	   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	   $(PROTO_FILE)

clean:
	rm $(OUT_FILE)
