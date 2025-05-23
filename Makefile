all: build/share/man/man1/tmux-tui.1

build:
	@mkdir -p build

build/tmux-tui build/docgen: build $(shell find . -type f -name "*.go")
	go build -o build ./...

build/share/man/man1/tmux-tui.1: build/docgen
	@build/docgen

clean:
	@rm -rf build
