package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/malparty/synth-xu/lib/generators"
	"github.com/malparty/synth-xu/lib/generators/effects"
	"github.com/malparty/synth-xu/lib/generators/oscillators"

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
	generator generators.Generator
}

func NewGame() *Game {
	return &Game{
		generator: *InitGenerators(),
	}
}

func (g *Game) Update() error {
	for i, key := range keys {
		if !inpututil.IsKeyJustPressed(key) {
			continue
		}

		g.generator.SetNote(i)
	}

	for _, key := range octavesKeys {
		if !inpututil.IsKeyJustPressed(key) {
			continue
		}

		if key == ebiten.KeyPeriod { // >
			g.generator.OctaveFreqUp()
		} else if key == ebiten.KeyComma { // >
			g.generator.OctaveFreqDown()
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

func InitGenerators() *generators.Generator {
	f := 512
	speaker.Init(beep.SampleRate(generators.SampleRate), 4800)
	// s, err := generators.SinTone(beep.SampleRate(48000), f)
	// if err != nil {
	// 	panic(err)
	// }

	limiter := &effects.Limiter{
		Rate: 20.0,
	}

	reverb := &effects.Reverb{
		MixRate:  100,
		FadeRate: 80,
		DelayMs:  50,
	}

	chainFunction := &generators.ChainGenerator{
		GeneratorFuncs: []generators.GeneratorFunction{
			oscillators.SawFunc,
			limiter.GetLimiterFunc(),
			reverb.GetReverbFunc(),
		},
	}

	s2, err := generators.NewGenerator(beep.SampleRate(48000), f, chainFunction.ChainFunc)
	if err != nil {
		panic(err)
	}
	// speaker.Play(s)
	speaker.Play(s2.GetOsc())

	return s2
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Piano (Ebiten Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
