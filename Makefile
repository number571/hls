GC=go build
GFILES=service.go initial.go config.go database.go models.go
.PHONY: default build clean
default: build
build: $(GFILES)
	$(GC) service.go initial.go config.go database.go models.go
clean:
	rm -f service service.db service.cfg
