.PHONY: all go rust

FFINDFILES := $(shell find . -type f -name 'ffind')

all: go rust c

clean:
	rm -fr $(FFINDFILES)

go:
	cd go ; go build

rust:
	cd rust ; cargo build

c:
	cd c ; $(CC) ffind.c -o ffind
