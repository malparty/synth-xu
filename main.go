// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build example
// +build example

package main

import (
	"fmt"
	"image/color"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/malparty/synth-xu/osc"
)

var (
	arcadeFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const (
		arcadeFontSize = 8
		dpi            = 72
	)
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    arcadeFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

const (
	screenWidth  = 320
	screenHeight = 240
)

// toBytes returns the 2ch little endian 16bit byte sequence with the given left/right sequence.
func toBytes(l, r []int16) []byte {
	if len(l) != len(r) {
		panic("len(l) must equal to len(r)")
	}
	b := make([]byte, len(l)*4)
	for i := range l {
		b[4*i] = byte(l[i])
		b[4*i+1] = byte(l[i] >> 8)
		b[4*i+2] = byte(r[i])
		b[4*i+3] = byte(r[i] >> 8)
	}
	return b
}

var (
	pianoImage = ebiten.NewImage(screenWidth, screenHeight)
)

func init() {
	const (
		keyWidth = 24
		y        = 48
	)

	whiteKeys := []string{"A", "S", "D", "F", "G", "H", "J", "K", "L"}
	for i, k := range whiteKeys {
		x := i*keyWidth + 36
		height := 112
		ebitenutil.DrawRect(pianoImage, float64(x), float64(y), float64(keyWidth-1), float64(height), color.White)
		text.Draw(pianoImage, k, arcadeFont, x+8, y+height-8, color.Black)
	}

	blackKeys := []string{"Q", "W", "", "R", "T", "", "U", "I", "O"}
	for i, k := range blackKeys {
		if k == "" {
			continue
		}
		x := i*keyWidth + 24
		height := 64
		ebitenutil.DrawRect(pianoImage, float64(x), float64(y), float64(keyWidth-1), float64(height), color.Black)
		text.Draw(pianoImage, k, arcadeFont, x+8, y+height-8, color.White)
	}

	octaveKeys := []string{"<", ">"}
	for i, k := range octaveKeys {
		x := i*keyWidth + 24
		y := 256
		height := 32
		ebitenutil.DrawRect(pianoImage, float64(x), float64(y), float64(keyWidth-1), float64(height), color.White)
		text.Draw(pianoImage, k, arcadeFont, x+8, y+height-8, color.White)
	}
}

var (
	keys = []ebiten.Key{
		ebiten.KeyQ,
		ebiten.KeyA,
		ebiten.KeyW,
		ebiten.KeyS,
		ebiten.KeyD,
		ebiten.KeyR,
		ebiten.KeyF,
		ebiten.KeyT,
		ebiten.KeyG,
		ebiten.KeyH,
		ebiten.KeyU,
		ebiten.KeyJ,
		ebiten.KeyI,
		ebiten.KeyK,
		ebiten.KeyO,
		ebiten.KeyL,
	}
	octavesKeys = []ebiten.Key{
		ebiten.KeyComma,
		ebiten.KeyPeriod,
	}
)

type Game struct {
	Osc *osc.Osc

	audioContext *audio.Context
	player       *audio.Player
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) InitOsc() error {
	if g.audioContext == nil {
		g.audioContext = audio.NewContext(osc.SampleRate)
	}

	if g.Osc == nil {
		g.Osc = osc.NewOsc()
	}

	if g.player == nil {
		// Pass the (infinite) stream to NewPlayer.
		// After calling Play, the stream never ends as long as the player object lives.
		var err error
		g.player, err = g.audioContext.NewPlayer(g.Osc)
		if err != nil {
			return err
		}
		g.player.Play()
	}
	return nil
}

func (g *Game) Update() error {
	g.InitOsc() // Starts oscillator stream

	for i, key := range keys {
		if !inpututil.IsKeyJustPressed(key) {
			continue
		}

		g.Osc.MoveNote(int64(i))
	}

	for _, key := range octavesKeys {
		if !inpututil.IsKeyJustPressed(key) {
			continue
		}

		if key == ebiten.KeyPeriod { // >
			g.Osc.MoveOctave(false)
		} else if key == ebiten.KeyComma { // >
			g.Osc.MoveOctave(true)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0x80, 0xc0, 0xff})
	screen.DrawImage(pianoImage, nil)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Piano (Ebiten Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
