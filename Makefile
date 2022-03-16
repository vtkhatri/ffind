.PHONY: all go rust

all: go rust c

go:
	cd go ; go build

rust:
	cd rust ; cargo build

c:
	cd c ; $(CC) ffind.c -o ffind
