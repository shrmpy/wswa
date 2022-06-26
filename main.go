// Copyright 2022 YOUREALLBREATHTAKING. All rights reserved.
// Use of this source code is subject to an MIT-style
// license which can be found in the LICENSE file.

//go:build js && wasm
// +build js,wasm

package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"syscall/js"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js dist/web/wasm_exec.js
//go:generate env GOOS=js GOARCH=wasm go build -v -ldflags "-w -s" -o dist/web/wswa.wasm ./

func main() {
	gameWidth, gameHeight := 640, 480

	ebiten.SetWindowSize(gameWidth, gameHeight)
	ebiten.SetWindowTitle("ebiten-websocket-wasm")

	game := &Game{
		Width:  gameWidth,
		Height: gameHeight,
		Player: &Player{image.Pt(gameWidth/2, gameHeight/2)},
		info:   make([]string, 0, 25),
	}

	game.jsws = wsconnect()
	defer game.jsws.Disconnect()
	game.jsws.Attach("message", game.wsread)
	rand.Seed(time.Now().UnixNano())

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

// Game represents the main game state
type Game struct {
	Width  int
	Height int
	Player *Player
	jsws   Socket
	info   []string
}

// Layout is hardcoded for now, may be made dynamic in future
func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return g.Width, g.Height
}

// Update calculates game logic
func (g *Game) Update() error {
	// Pressing Q any time quits immediately
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return errors.New("game quit by player")
	}

	// Pressing F toggles full-screen
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		} else {
			ebiten.SetFullscreen(true)
		}
	}

	// Movement controls
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Player.Move()

		pos := fmt.Sprintf("pos(%d, %d); roll: %d",
			g.Player.Coords.X, g.Player.Coords.Y,
			rand.Intn(g.Player.Coords.Y))
		g.jsws.Send(pos)
	}

	return nil
}

// Draw draws the game screen by one frame
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(
		screen,
		float64(g.Player.Coords.X),
		float64(g.Player.Coords.Y),
		20,
		20,
		color.White,
	)

	g.wsprint(screen)
}

// Player is the player character in the game
type Player struct {
	Coords image.Point
}

// Move moves the player upwards
func (p *Player) Move() {
	p.Coords.Y--
}

// print echoes from ws server
func (g *Game) wsprint(screen *ebiten.Image) {
	for i := len(g.info); i > 0; i-- {
		y := g.Height - i*lineHt
		ebitenutil.DebugPrintAt(screen, g.info[i-1], 0, y)
	}
}

// attached to message event from websock
func (g *Game) wsread(ev js.Value) error {
	// arbitrary 25 line display
	if len(g.info) > 25 {
		return nil
	}
	var data = ev.Get("data").String()
	g.info = append(g.info, data)

	return nil
}

// open connection to ws server
func wsconnect() Socket {
	var ws = js.Global().Get("WebSocket").New("ws://localhost:8077/")
	return Socket{Value: ws}
}

// send text to ws server
func (w Socket) Send(txt string) {
	w.Call("send", txt)
}

// register event handler
func (w Socket) Attach(ev string, fn func(js.Value) error) {
	var cb = js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			fn(args[0])
			return nil
		})
	w.handlers = append(w.handlers, cb)
	w.Call("addEventListener", ev, cb)
}

// disco from ws server
func (w Socket) Disconnect() {
	for _, cb := range w.handlers {
		// clean-up
		cb.Release()
	}
	w.Call("close")
}

// strong type for the Js obj
type Socket struct {
	handlers []js.Func
	js.Value
}

const lineHt = 14
