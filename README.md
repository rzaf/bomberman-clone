# bomberman-clone
Super bomberman 2 (SNES) clone in Go using <a href="https://github.com/gen2brain/raylib-go">raylib-go</a>

## trying out
download latest executable from <a href="https://github.com/rzaf/bomberman-clone/releases">releases</a>

![Main menu](screenshots/2.png)

![Battle](screenshots/1.png)

![Level editor](screenshots/3.png)

## building
* install protocol buffer compiler (<a href="https://grpc.io/docs/protoc-installation/">link</a>) if you want to compile pb files
* clone project and get get into directory `git clone https://github.com/rzaf/bomberman-clone.git && cd bomberman-clone`
* run `go mod download` to get required modules
* run `make all` or run `make build` if you dont have protoc installed
* run `bomberman-clone` in `bin`


## features
* Upgrades (speed up, wall pass, extra bomb, ...)
* Local multiplayer
* Online multiplayer (grpc) **WIP**
* Option menu (audio and keymapping)
* Controller support
* Level editor
