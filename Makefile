GC=go build
.PHONY: default build clean
default: build
build:
	$(GC) -o cmd/client/c-hls cmd/client/*.go
	$(GC) -o cmd/server/s-hls cmd/server/*.go
	$(GC) -o cmd/service/hls cmd/service/*.go
clean:
	rm -f cmd/client/c-hls cmd/server/s-hls cmd/service/hls 
