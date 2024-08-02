ifeq ($(OS),Windows_NT) 
    DETECTED_OS := Windows
	FILE := bomberman-clone.exe
	PB_FILE := pb\game.proto
else
    DETECTED_OS := $(shell uname)
	FILE := bomberman-clone
	PB_FILE := pb/game.proto
endif

all: buildPbs build

buildPbs:
	@echo "generating protocol buffer files ... "
	@protoc ${PB_FILE} --go_out=. --go-grpc_out=.

build:
ifeq ($(DETECTED_OS),Windows)
	@if not exist bin mkdir bin
	@echo "compiling go files (may takes too long for first time)" 
	@set CGO_ENABLED=1 && set CC=gcc && go build -ldflags "-H=windowsgui" -o bin\${FILE} main.go
	@echo "moving assets to bin"
	@if not exist bin\assets mkdir bin\assets
	@xcopy /I /Y /Q assets\audio bin\assets\audio
	@xcopy /I /Y /Q assets\maps bin\assets\maps
	@copy /Y assets\menuMap.txt bin\assets
	@copy /Y assets\characters.png bin\assets
	@copy /Y assets\tiles.png bin\assets
else
	@mkdir -p bin
	@echo "compiling go files (may takes too long for first time)" 
	@CGO_ENABLED=1 go build -o bin/${FILE} main.go
	@echo "moving assets to bin"
	@mkdir -p bin/assets
	@cp -r assets/audio bin/assets/
	@cp -r assets/maps bin/assets/
	@cp assets/menuMap.txt bin/assets/
	@cp assets/characters.png bin/assets/
	@cp assets/tiles.png bin/assets/
endif

clean:
	@echo "Removing bin directory ... "
ifeq ($(DETECTED_OS),Windows)
	@if exist bin rmdir /s /q bin 
else
	@rm -rf bin 
endif