package main

import (
	"fmt"
	"math"
	"math/cmplx"

	rl "github.com/gen2brain/raylib-go/raylib"
	"gonum.org/v1/gonum/dsp/fourier"
)

const (
	screenWidth  = 800
	screenHeight = 600

	bufSize = 4096
)

var (
	audioBuf = make([]float64, bufSize, 2*bufSize)

	freqBuf = make([]complex128, bufSize/2+1)
	fft     = fourier.NewFFT(bufSize)

	musLoaded = false
	mus       rl.Music
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Music Visualizer")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.AttachAudioMixedProcessor(streamProcessor)
	defer rl.DetachAudioMixedProcessor(streamProcessor)

	visRect := rl.Rectangle{
		X:      100.0,
		Y:      100.0,
		Width:  float32(rl.GetScreenWidth()) - 200.0,
		Height: float32(rl.GetScreenHeight()) - 220.0,
	}

	for !rl.WindowShouldClose() {
		update()

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		if musLoaded {
			drawMusicBar()
			drawVis(audioBuf, visRect)
		} else {
			fontSize := int32(24)
			text := "Drop a music file here"
			textWidth := rl.MeasureText(text, fontSize)
			rl.DrawText(text, (int32(rl.GetScreenWidth())-textWidth)/2, int32(rl.GetScreenHeight())/2-12, fontSize, rl.White)
		}

		rl.EndDrawing()
	}

	if musLoaded {
		rl.UnloadMusicStream(mus)
	}
}

func update() {
	if rl.IsFileDropped() {
		droppedFiles := rl.LoadDroppedFiles()
		if len(droppedFiles) > 0 {
			if musLoaded {
				rl.UnloadMusicStream(mus)
			}
			mus = rl.LoadMusicStream(droppedFiles[0])
			rl.PlayMusicStream(mus)
			musLoaded = true
		}
	}

	if !musLoaded {
		return
	}

	rl.UpdateMusicStream(mus)

	if rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyP) {
		if rl.IsMusicStreamPlaying(mus) {
			rl.PauseMusicStream(mus)
		} else {
			rl.ResumeMusicStream(mus)
		}
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mousePos := rl.GetMousePosition()
		if mousePos.Y > float32(rl.GetScreenHeight())-20.0 {
			seekPos := mousePos.X / float32(rl.GetScreenWidth())
			if seekPos >= 0 && seekPos <= 1 {
				rl.SeekMusicStream(mus, seekPos*rl.GetMusicTimeLength(mus))
			}
		}
	}
}

func streamProcessor(data []float32, frames int) {
	for i := 0; i < len(data); i += 2 {
		var sl float64 = float64(data[i])
		var sr float64
		if i+1 < len(data) {
			sr = float64(data[i+1])
		}
		s := math.Max(sl, sr)
		audioBuf = append(audioBuf, s)
	}
	if len(audioBuf) > bufSize {
		audioBuf = audioBuf[len(audioBuf)-bufSize:]
	}
}

func drawVis(data []float64, bounds rl.Rectangle) {
	freqBuf = fft.Coefficients(freqBuf, audioBuf)

	barWidth := float32(bounds.Width) / float32((len(freqBuf)))

	for i, f := range freqBuf {
		barHeight := float32(math.Min(cmplx.Abs(f), 1.0)) * bounds.Height
		rl.DrawRectangleV(
			rl.Vector2{
				X: float32(i)*(barWidth) + bounds.X,
				Y: bounds.Y + bounds.Height - barHeight,
			},
			rl.Vector2{
				X: barWidth,
				Y: barHeight,
			},
			rl.Green,
		)
	}
}

func drawMusicBar() {
	timePlayed := rl.GetMusicTimePlayed(mus)
	musLength := rl.GetMusicTimeLength(mus)

	sw := int32(rl.GetScreenWidth())
	sh := int32(rl.GetScreenHeight())

	var barHeight int32 = 20
	rl.DrawRectangle(0, sh-barHeight, sw, barHeight, rl.Gray)
	rl.DrawRectangle(0, sh-barHeight, int32(timePlayed/musLength*float32(sw)), barHeight, rl.Maroon)

	text := fmt.Sprintf("%.0f / %.0f", timePlayed, musLength)
	var fs int32 = 14
	tw := rl.MeasureText(text, fs)
	rl.DrawText(text, sw/2-tw/2, sh-barHeight+3, 14, rl.White)
}
