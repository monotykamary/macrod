.PHONY: build build-daemon build-tui run-daemon run-tui clean

build: build-daemon build-tui

build-daemon:
	go build -o bin/macrod-daemon cmd/daemon/main.go

build-tui:
	go build -o bin/macrod-tui cmd/tui/main.go

run-daemon:
	./bin/macrod-daemon

run-tui:
	./bin/macrod-tui

clean:
	rm -rf bin/

install-deps:
	go get github.com/charmbracelet/bubbletea
	go get github.com/charmbracelet/bubbles
	go get github.com/charmbracelet/lipgloss
	go get github.com/robotgo/gohook