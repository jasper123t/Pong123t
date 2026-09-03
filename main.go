package main

import (
	"fmt"
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
	BallRadius = 15

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
	b.x += math.Sin(b.dir*(math.Pi/180.0)) * 10
	b.y += math.Cos(b.dir*(math.Pi/180.0)) * 10
	ctx.Set("fillStyle", BallColor)
	ctx.Call("beginPath")
	ctx.Call("arc", b.x, b.y, BallRadius, 0, 2*math.Pi)
	ctx.Call("fill")
}

func (b *Ball) bounce(delta float64) {
	b.dir = delta - b.dir
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
		dir: float64(randInt(160, 200)),
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

		resx := gameBall.x + math.Sin(gameBall.dir*(math.Pi/180.0))*10 + BallRadius
		resy := gameBall.y + math.Cos(gameBall.dir*(math.Pi/180.0))*10
		resd := ctx.Call("getImageData", resx, resy, 1, 1).Get("data")
		// 	resd.Index(0).Int(), // R
		// 	resd.Index(1).Int(), // G
		// 	resd.Index(2).Int(), // B

		if keysPressed["KeyA"] || keysPressed["ArrowLeft"] {
			gameBar.x -= 5
		}
		if keysPressed["KeyD"] || keysPressed["ArrowRight"] {
			gameBar.x += 5
		}

		if gameBall.x < BallRadius {
			gameBall.bounce(0)
			gameBall.x = BallRadius
		}
		if gameBall.x > CanvasWidth-BallRadius {
			gameBall.bounce(0)
			gameBall.x = CanvasWidth - BallRadius
		}
		if gameBall.y < BallRadius {
			gameBall.bounce(180)
			gameBall.y = BallRadius
		}
		if resd.Index(1).Int() == 153 {
			gameBall.bounce(gameBall.dir + float64(180+randInt(-20, 20)+randInt(-20, 20)))
		}
		if gameBall.y > CanvasHeight-31 {
			fmt.Printf("dead\n")
		}
		if gameBall.y > CanvasHeight-BallRadius {
			gameBall.bounce(180)
			gameBall.y = CanvasHeight - BallRadius
		}
		ctx.Call("clearRect", 0, 0, CanvasWidth, CanvasHeight)
		gameBar.Draw(ctx)
		gameFloor.Draw(ctx)
		gameBall.Draw(ctx)
	}

	select {}

}
