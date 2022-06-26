# ebiten-websocket-wasm

Take [ebiten-game-template](https://github.com/sinisterstuf/ebiten-game-template)
 and add websockets. Will it work?

Following the goita project, it reveals more clearly how there are two layers
 one for the service, and a second for the webview. The service can utilize packages
 written for websockets (i.e., `gorilla/websocket`). The webview where the WASM runtime
 lives, must rely on `syscall/js` to "talk" to the backend.

-----✂️-----

> ⚠️ After cloning this repository:

> Write your OWN name name in the LICENSE file and run this command to replace the game name (tested on Linux and Mac):

```bash
grep -Rl ebiten-game-template | xargs sed -i '' -e "s/ebiten-game-template/${PWD##*/}/g"
```

> it assumes that the game name is the name of the current folder because that is what `go build` will call it.

> Then delete this section from the README, and start editing `main.go` to make your own game!

-----✂️-----

## For game testers

<!-- TODO: add a link to the latest downloads page -->

Game controls:
- F: toggle full-screen
- Q: quit the game
- Space: move up

## For programmers

Make sure you have [Go 1.17 or later](https://go.dev/) to contribute to the game

To build the game yourself, run: `go build .` it will produce an ebiten-game-template file and on Windows ebiten-game-template.exe.

To run the tests, run: `go test ./...` but there are no tests yet.

The project has a very simple, flat structure, the first place to start looking is the main.go file.

## Credits

Websocket example is
 by [Sergey Kamardin](https://github.com/gobwas/ws)
 ([LICENSE](https://github.com/gobwas/ws/blob/master/LICENSE))

Websocket + WASM + webview
 by [Marc](https://github.com/Markcial/goita)
 ([LICENSE](https://github.com/Markcial/goita/blob/main/LICENSE))

Github workflow
 by [Siôn le Roux](https://github.com/sinisterstuf/ebiten-game-template)
 ([LICENSE](https://github.com/sinisterstuf/ebiten-game-template/blob/main/LICENSE))

Ebitengine
 by [Hajime Hoshi](https://github.com/hajimehoshi/ebiten/)
 ([LICENSE](https://github.com/hajimehoshi/ebiten/blob/main/LICENSE))

