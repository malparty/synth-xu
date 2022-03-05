package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/malparty/synth-xu/lib/constant"
	"github.com/malparty/synth-xu/lib/modules"
	"github.com/malparty/synth-xu/lib/modules/effects"
	"github.com/malparty/synth-xu/lib/modules/modulators"
	"github.com/malparty/synth-xu/lib/modules/oscillators"
	"github.com/malparty/synth-xu/lib/racks"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
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
		y := 175
		height := 32
		ebitenutil.DrawRect(pianoImage, float64(x), float64(y), float64(keyWidth-1), float64(height), color.White)
		text.Draw(pianoImage, k, arcadeFont, x+8, y+height-8, color.Black)
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
	voice    *racks.Voice
	envelope *modulators.Envelope
	osc      *oscillators.MultiOsc
}

func NewGame() *Game {
	game := &Game{}

	game.initAudioModules()

	return game
}

func (g *Game) Update() error {
	for i, key := range keys {
		if inpututil.IsKeyJustPressed(key) {
			g.voice.SetNote(i)
			g.envelope.TriggerNote()
		} else if inpututil.IsKeyJustReleased(key) {
			g.envelope.ReleaseNote()
		}
	}

	for _, key := range octavesKeys {
		if !inpututil.IsKeyJustPressed(key) {
			continue
		}

		if key == ebiten.KeyPeriod { // >
			g.voice.OctaveFreqUp()
		} else if key == ebiten.KeyComma { // >
			g.voice.OctaveFreqDown()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		g.osc.OscA.Type = oscillators.Saw
		g.displayCurrentChainCycle()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.osc.OscA.Type = oscillators.Sin
		g.displayCurrentChainCycle()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.osc.OscA.Type = oscillators.Triangle
		g.displayCurrentChainCycle()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		g.osc.OscB.Type = oscillators.Square
		g.displayCurrentChainCycle()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		g.osc.OscB.Type = oscillators.Saw
		g.displayCurrentChainCycle()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.osc.OscB.Type = oscillators.Sin
		g.displayCurrentChainCycle()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if g.osc.MixPercent < 96 {
			g.osc.MixPercent += 5
			g.displayCurrentChainCycle()
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if g.osc.MixPercent > 4 {
			g.osc.MixPercent -= 5
			g.displayCurrentChainCycle()
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

func (g *Game) initAudioModules() {
	f := 512
	speaker.Init(beep.SampleRate(constant.SampleRate), 4800)

	oscA := &oscillators.Osc{
		Type: oscillators.Saw,
	}
	oscB := &oscillators.Osc{
		Type: oscillators.Sin,
	}
	g.osc = &oscillators.MultiOsc{
		OscA:       oscA,
		OscB:       oscB,
		MixPercent: 0,
	}

	limiter := &effects.Limiter{
		Rate: 20.0,
	}

	// Improvement: Make a rack lane for effects  and apply it AFTER the ADSR envelope!
	// reverb := &effects.Reverb{
	// 	MixRate:  100,
	// 	FadeRate: 50,
	// 	DelayMs:  40,
	// }

	g.envelope = &modulators.Envelope{
		Attack:  0.1,
		Decay:   0.1,
		Sustain: 0.8,
		Release: 0.3,
	}

	chainFunction := racks.NewChainFunc(g.envelope, []modules.Module{
		g.osc,
		limiter,
		// reverb,
	},
	)

	voice, err := racks.NewVoice(beep.SampleRate(48000), f, chainFunction)
	if err != nil {
		panic(err)
	}
	g.voice = voice

	speaker.Play(voice.GetOsc())

	g.displayCurrentChainCycle()

}

func (g *Game) displayCurrentChainCycle() {
	positionX := 100.0
	positionY := 200.0

	previousX := 0.0
	previousY := 0.0

	delta := 0.01

	// Remove existing content on this area:
	ebitenutil.DrawRect(pianoImage, positionX, positionY-40, 200, 80, color.RGBA{
		R: 200,
		G: 100,
		B: 200,
		A: 255,
	})

	// Build a sample to dysplay on the screen
	g.envelope.TriggerNote()

	for i := 0.0; i < 2; i += delta {
		_, stat := math.Modf(i)
		stat = g.voice.ChainFunction.ChainFunc(stat, delta)

		x := i * 100
		y := stat * 100
		ebitenutil.DrawLine(pianoImage, positionX+previousX, positionY+previousY, positionX+x, positionY+y, color.Black)

		previousX = x
		previousY = y
	}

	g.envelope.ReleaseNote()
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Piano (Ebiten Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
