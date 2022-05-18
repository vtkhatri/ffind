.PHONY: all build clean go rust

FFINDFILES := $(shell find . -type f -name 'ffind')

all: build
build: gobuild rustbuild cbuild

clean:
	rm -fr $(FFINDFILES)

gobuild:
	cd go && go build

go: gobuild
	cd go && go install

rustbuild:
	cd rust && cargo build

rust: rustbuild
	cd rust && cargo install

cbuild:
	cd c && $(CC) -g ffind.c -o ffind
