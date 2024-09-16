
BINARY=deployme
GOFILES=deployme.go

# Compiler settings
GOOS_LINUX=linux
GOOS_WINDOWS=windows
GOARCH=amd64

.PHONY: all clean windows linux

# Default rule - compile for both platforms
both: windows linux

# Compile for Windows
windows:
	GOOS=$(GOOS_WINDOWS) GOARCH=$(GOARCH) go build -o bin/$(BINARY)-windows.exe $(GOFILES)

# Compile for Linux
linux:
	GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH) go build -o bin/$(BINARY)-linux $(GOFILES)

# Clean up compiled binaries
clean:
	rm -f $(BINARY)-windows.exe $(BINARY)-linux

run:
	go run deployme.go