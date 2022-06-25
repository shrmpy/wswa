// Copyright 2022 YOUREALLBREATHTAKING.  All rights reserved.
// Use of this source code is subject to an MIT-style
// license which can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/gobwas/ws/wsutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func main() {
	gameWidth, gameHeight := 640, 480

	ebiten.SetWindowSize(gameWidth, gameHeight)
	ebiten.SetWindowTitle("ebiten-game-template")

	game := &Game{
		Width:  gameWidth,
		Height: gameHeight,
		Player: &Player{image.Pt(gameWidth/2, gameHeight/2)},
		info:   make([]string, 0, 25),
	}

	game.wsc = wsconnect()
	defer game.wsc.Close()
	rand.Seed(time.Now().UnixNano())

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Game represents the main game state
type Game struct {
	Width  int
	Height int
	Player *Player
	wsc    net.Conn
	info   []string
}

// Layout is hardcoded for now, may be made dynamic in future
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Width, g.Height
}

// Update calculates game logic
func (g *Game) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go g.wsread(ctx)

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
		g.wsrequest(pos)
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

// send text to the ws server
func (g *Game) wsrequest(msg string) {
	wsutil.WriteClientText(g.wsc, []byte(msg))
}

// print echoes from ws server
func (g *Game) wsprint(screen *ebiten.Image) {
	for i := len(g.info); i > 0; i-- {
		y := g.Height - i*lineHt
		ebitenutil.DebugPrintAt(screen, g.info[i-1], 0, y)
	}
}

// read message from websock
func (g *Game) wsread(ctx context.Context) error {
	// non-polling, limit waiting (with ctx)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// arbitrary 25 line display
		if len(g.info) > 25 {
			return nil
		}
		buf, err := wsutil.ReadServerText(g.wsc)
		if err != nil {
			return err
		}
		g.info = append(g.info, string(buf))
	}
	return nil
}

// open connection to ws server
func wsconnect() net.Conn {
	var dd = wsutil.DebugDialer{OnResponse: debugrcv, OnRequest: debugreq}
	conn, _, _, err := dd.Dial(context.Background(), "ws://localhost:8077/")
	if err != nil {
		log.Fatalf("FAIL websock dial, %s", err.Error())
	}
	return conn
}
func debugreq(data []byte) {
	log.Printf("REQ: ", string(data))
}
func debugrcv(data []byte) {
	log.Printf("RCV: ", string(data))
}

const lineHt = 14
