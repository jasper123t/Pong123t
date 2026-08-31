package main

import (
	"math"
	"math/rand/v2"
	"syscall/js"
	"time"
)

const (
	FPS = time.Second / 30

	CanvasWidth  = 480
	CanvasHeight = 360

	BallColor  = "#ff0000"
	BallRadius = "15"

	BarColor  = "#009933"
	BarWidth  = 186
	BarHeight = 10

	FloorColor  = "#003fff"
	FloorWidth  = 478
	FloorHeight = 30
)

type Ball struct {
	x   float64
	y   float64
	dir float64
}

func (b *Ball) Draw(ctx js.Value) {
	b.y += math.Sin(b.dir*(math.Pi/180.0)) * 10
	b.x += math.Cos(b.dir*(math.Pi/180.0)) * 10
	ctx.Set("fillStyle", BallColor)
	ctx.Call("beginPath")
	ctx.Call("arc", b.x, b.y, BallRadius, 0, 2*math.Pi)
	ctx.Call("fill")
}

type Bar struct {
	x int
}

func (b *Bar) Draw(ctx js.Value) {
	ctx.Set("fillStyle", BarColor)
	ctx.Call("beginPath")
	ctx.Call("arc", b.x-BarWidth/2, 292+BarHeight/2, BarHeight/2, 0, 2*math.Pi)
	ctx.Call("rect", b.x-BarWidth/2, 292, BarWidth, BarHeight)
	ctx.Call("arc", b.x+BarWidth/2, 292+BarHeight/2, BarHeight/2, 0, 2*math.Pi)
	ctx.Call("fill")
}

type Floor struct{}

func (f *Floor) Draw(ctx js.Value) {
	ctx.Set("fillStyle", FloorColor)
	ctx.Call("beginPath")
	ctx.Call("rect", 0, CanvasHeight-31, FloorWidth, FloorHeight)
	ctx.Call("fill")
}

func randInt(min, max int) int {
	return rand.IntN(max-min+1) + min
}

func main() {
	window := js.Global()
	document := window.Get("document")
	canvas := document.Call("getElementById", "gameCanvas")
	ctx := canvas.Call("getContext", "2d")

	gameBall := Ball{
		x:   float64(CanvasWidth/2 + randInt(-10, 10)),
		y:   float64(CanvasHeight/2 + randInt(-10, 10)),
		dir: float64(randInt(70, 110)),
	}

	gameBar := Bar{
		x: CanvasWidth / 2,
	}

	gameFloor := Floor{}

	var keysPressed = make(map[string]bool)

	onKeyDown := js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		key := event.Get("code").String()
		keysPressed[key] = true
		return nil
	})

	onKeyUp := js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		key := event.Get("code").String()
		keysPressed[key] = false
		return nil
	})

	js.Global().Call("addEventListener", "keydown", onKeyDown)
	js.Global().Call("addEventListener", "keyup", onKeyUp)

	for range time.Tick(FPS) {
		if keysPressed["KeyA"] || keysPressed["ArrowLeft"] {
			gameBar.x -= 5
		}
		if keysPressed["KeyD"] || keysPressed["ArrowRight"] {
			gameBar.x += 5
		}
		ctx.Call("clearRect", 0, 0, CanvasWidth, CanvasHeight)
		gameBall.Draw(ctx)
		gameBar.Draw(ctx)
		gameFloor.Draw(ctx)
	}

	select {}

}
