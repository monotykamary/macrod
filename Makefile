.PHONY: build build-daemon build-tui run-daemon run-tui clean release release-darwin-amd64 release-darwin-arm64

build: build-daemon build-tui

build-daemon:
	CGO_ENABLED=1 go build -o bin/macrod-daemon cmd/daemon/main.go

build-tui:
	go build -o bin/macrod-tui cmd/tui/main.go

run-daemon:
	./bin/macrod-daemon

run-tui:
	./bin/macrod-tui

clean:
	rm -rf bin/ dist/

install-deps:
	go get github.com/charmbracelet/bubbletea
	go get github.com/charmbracelet/bubbles
	go get github.com/charmbracelet/lipgloss
	go get github.com/robotgo/gohook

# Release targets
release-darwin-amd64:
	mkdir -p dist
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o dist/macrod-daemon cmd/daemon/main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/macrod-tui cmd/tui/main.go
	cd dist && tar -czf macrod-darwin-amd64.tar.gz macrod-daemon macrod-tui

release-darwin-arm64:
	mkdir -p dist
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o dist/macrod-daemon cmd/daemon/main.go
	GOOS=darwin GOARCH=arm64 go build -o dist/macrod-tui cmd/tui/main.go
	cd dist && tar -czf macrod-darwin-arm64.tar.gz macrod-daemon macrod-tui

release: clean release-darwin-amd64 release-darwin-arm64
	@echo "Release archives created in dist/"
	@echo "Upload these to GitHub releases:"
	@echo "  - dist/macrod-darwin-amd64.tar.gz"
	@echo "  - dist/macrod-darwin-arm64.tar.gz"