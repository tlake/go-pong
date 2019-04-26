package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth    int = 800
	winHeight   int = 600
	paddleSpeed     = 400
	ballSpeed       = 400
)

var gameIsActive = false
var gameScores = map[string]int{
	"leftPaddle":  0,
	"rightPaddle": 0,
}

var nums = [][]byte{
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		0, 1, 0,
		1, 0, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		0, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	{
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		0, 0, 1,
	},
	{
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
}

func drawNumber(pos pos, color color, size, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

func getCenter() pos {
	return pos{float32(winWidth) / 2, float32(winHeight) / 2}
}

type pos struct {
	x, y float32
}

type color struct {
	r, g, b byte
}

type ball struct {
	pos
	radius float32
	xv     float32
	yv     float32
	color  color
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

func (ball *ball) update(leftPaddle, rightPaddle *paddle, elapsedTime float32, gameIsActive *bool) {
	ball.x += ball.xv * elapsedTime
	ball.y += ball.yv * elapsedTime

	// handle collisions
	//	top and bottom of screen
	if ball.y-ball.radius < 0 {
		ball.y = ball.radius
		ball.yv = -ball.yv
	}

	if ball.y+ball.radius > float32(winHeight) {
		ball.y = float32(winHeight) - ball.radius
		ball.yv = -ball.yv
	}

	//	left and right of screen
	if ball.x < 0 {
		gameScores["rightPaddle"]++
		ball.pos = getCenter()
		*gameIsActive = false
	}

	if ball.x > float32(winWidth) {
		gameScores["leftPaddle"]++
		ball.pos = getCenter()
		*gameIsActive = false
	}

	//	paddles
	if ball.x-ball.radius < leftPaddle.x+leftPaddle.w/2 && ball.x-ball.radius > leftPaddle.x-leftPaddle.w/2 {
		if ball.y > leftPaddle.y-leftPaddle.h/2 && ball.y < leftPaddle.y+leftPaddle.h/2 {
			ball.x = leftPaddle.x + leftPaddle.h/2
			ball.xv = -ball.xv
		}
	}

	if ball.x+ball.radius > rightPaddle.x-rightPaddle.w/2 && ball.x+ball.radius < rightPaddle.x+rightPaddle.w/2 {
		if ball.y > rightPaddle.y-rightPaddle.h/2 && ball.y < rightPaddle.y+rightPaddle.h/2 {
			ball.x = rightPaddle.x - rightPaddle.h/2
			ball.xv = -ball.xv
		}
	}
}

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	color color
}

func (paddle *paddle) draw(pixels []byte) {
	startX := paddle.x - paddle.w/2
	startY := paddle.y - paddle.h/2

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixel(int(startX)+x, int(startY)+y, paddle.color, pixels)
		}
	}
}

func (paddle *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		if paddle.y-paddle.h/2 > 0 {
			paddle.y -= paddle.speed * elapsedTime
		}
	}

	if keyState[sdl.SCANCODE_DOWN] != 0 {
		if paddle.y+paddle.h/2 < float32(winHeight) {
			paddle.y += paddle.speed * elapsedTime
		}
	}
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	if ball.y > paddle.h/2 && ball.y < float32(winHeight)-paddle.h/2 {
		paddle.y = ball.y
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	idx := (y*int(winWidth) + x) * 4

	if idx < len(pixels)-4 && idx > -0 {
		pixels[idx] = c.r
		pixels[idx+1] = c.g
		pixels[idx+2] = c.b
	}
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Go Pong",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth),
		int32(winHeight),
		sdl.WINDOW_SHOWN,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos{100, 300}, 20, 100, paddleSpeed, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth) - 100, 300}, 20, 100, paddleSpeed, color{255, 255, 255}}
	ball := ball{getCenter(), 20, ballSpeed, ballSpeed, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	var frameStart time.Time
	var elapsedTime float32

	tex.Update(nil, pixels, int(winWidth)*4)
	renderer.Copy(tex, nil, nil)
	renderer.Present()

	// GAME LOOP
	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		clear(pixels)

		drawNumber(pos{200, 75}, color{255, 255, 255}, 20, gameScores["leftPaddle"], pixels)
		drawNumber(pos{float32(winWidth - 200), 75}, color{255, 255, 255}, 20, gameScores["rightPaddle"], pixels)

		player1.update(keyState, elapsedTime)
		player2.aiUpdate(&ball, elapsedTime)

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		if gameIsActive {
			ball.update(&player1, &player2, elapsedTime, &gameIsActive)
			ball.draw(pixels)
		} else {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				gameIsActive = true
			}
		}

		tex.Update(nil, pixels, int(winWidth)*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < .005 {
			sdl.Delay(5 - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}

}
