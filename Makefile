ifeq ($(OS),Windows_NT) 
    detected_OS := Windows
	file := bomberman-clone.exe
else
    detected_OS := $(shell uname)
	file := bomberman-clone
endif

all: buildPbs build

buildPbs:
	@echo "generating protocol buffer files ... "
	@protoc pb/game.proto --go_out=. --go-grpc_out=.

build:
	@mkdir -p bin
	@echo "compiling go files (may takes too long for first time)" 
	@go build -o bin/${file} main.go
	@echo "moving assets to bin"
	@mkdir -p bin/assets
	@cp -r assets/audio bin/assets/
	@cp -r assets/maps bin/assets/
	@cp assets/menuMap.txt bin/assets/
	@cp assets/characters.png bin/assets/
	@cp assets/tiles.png bin/assets/

clean:
	@echo "Removing bin directory ... "
	@rm -rf bin/