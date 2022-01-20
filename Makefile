.PHONY: all go rust

all: go rust

go:
	cd go ; go build

rust:
	cd rust ; cargo build
