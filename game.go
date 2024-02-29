package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/bitmapfont/v3"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	mode        mode
	msg         string
	cnt         int
	score       int
	highscore   int
	hitScore    int
	maxHitScore int
	bird        bird
	balls       balls
}

type mode int

const (
	title mode = iota
	playing
	gameover
)

type bird struct {
	images map[bird_Img_Structure]*ebiten.Image
	size   float64
	color  string
	state  string
	x      float64
	y      float64
	speed  float64
}

type balls struct {
	images        map[string]*ebiten.Image
	number        int
	ball          []ball
	accelerations []float64
}

type ball struct {
	size         float64
	color        string
	x            float64
	y            float64
	speed        float64
	acceleration float64
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() error {
	g.cnt = 0
	g.score = 0
	g.hitScore = 0
	g.birdInit()
	g.ballsInit()
	return nil
}

func (g *Game) birdInit() error {
	// g.bird.imagesは変数birdImagesを参照する
	g.bird.images = birdImages
	g.bird.size = birdSize
	g.bird.color = "red"
	g.bird.state = "up"
	g.bird.x = float64(screenWidth)/2 - birdSize/2
	g.bird.y = 240
	g.bird.speed = 2.5
	return nil
}

func (g *Game) ballsInit() error {
	// g.balls.imagesは変数ballImagesを参照する
	g.balls.images = ballImages
	g.balls.number = 1
	g.balls.accelerations = []float64{0.05}
	g.balls.ball = nil
	s := ball{}
	g.balls.ball = append(g.balls.ball, s)
	for i := range g.balls.ball {
		g.ballInit(i)
	}
	return nil
}

func (g *Game) ballInit(index int) error {
	g.balls.ball[index].size = ballSize
	g.balls.ball[index].color = colors[rand.IntN(len(colors[0:3]))]
	g.balls.ball[index].x = float64(rand.IntN(screenWidth - int(g.balls.ball[index].size)))
	g.balls.ball[index].y = -ballSize
	g.balls.ball[index].speed = 1
	g.balls.ball[index].acceleration = g.balls.accelerations[rand.IntN(len(g.balls.accelerations))]
	return nil
}

func (g *Game) Update() error {
	switch g.mode {
	case title:
		g.msg = fmt.Sprint("title, start:space key")
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.init()
			g.mode = playing
		}
	case playing:
		g.msg = fmt.Sprintf("score:%d, highscore:%d, hit:%d, maxhit:%d, birdspeed:%.2f, now:playing", g.score, g.highscore, g.hitScore, g.maxHitScore, g.bird.speed)
		g.updateBird()
		g.updateBalls()
		g.updateBirdAndBalls()
		g.cnt++
		if g.cnt%10 == 0 {
			g.score++
		}
	case gameover:
		g.msg = fmt.Sprintf("score:%d, highscore:%d, hit:%d, maxhit:%d, birdspeed:%.2f, now:gameover, restart:space key", g.score, g.highscore, g.hitScore, g.maxHitScore, g.bird.speed)
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.init()
			g.mode = playing
		}
	}
	return nil
}

func (g *Game) updateBird() {
	g.birdAnimation()
	g.birdMove()
	g.changeBirdColor()
	g.processOutScreen_Bird()
}

func (g *Game) updateBalls() {
	g.ballsMove()
	g.processOutScreen_Balls()
	g.addBall()
	g.increaseBallAcceleration()
}

func (g *Game) updateBirdAndBalls() {
	g.processHitBall()
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 背景
	g.drawBackground(screen)
	// テキスト
	text.Draw(screen, g.msg, bitmapfont.Face, 5, 10, color.White)
	text.Draw(screen, "change color key:[next color:'c', Red:'1', Yellow:'2', Green:'3'], move key:[Left:'<-', Right:'->']", bitmapfont.Face, 5, 476, color.White)
	// ボール
	g.drawBalls(screen)
	// 鳥
	g.drawBird(screen)
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(backgroundImage, op)
}

func (g *Game) drawBalls(screen *ebiten.Image) {
	for i := range g.balls.ball {
		g.drawBall(screen, i)
	}
}

func (g *Game) drawBall(screen *ebiten.Image, index int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.balls.ball[index].x, g.balls.ball[index].y)
	screen.DrawImage(g.balls.images[g.balls.ball[index].color], op)
}

func (g *Game) drawBird(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.bird.x, g.bird.y)
	screen.DrawImage(g.bird.images[bird_Img_Structure{g.bird.color, g.bird.state}], op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// 鳥を動かす
func (g *Game) birdMove() {
	// 右キー
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.bird.x <= 555 {
			g.bird.x += g.bird.speed
		}
	}
	// 左キー
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if g.bird.x >= 0 {
			g.bird.x -= g.bird.speed
		}
	}
}

// 鳥の色を変える
func (g *Game) changeBirdColor() {
	// 1キー:赤、 2キー:青、3キー:緑
	if ebiten.IsKeyPressed(ebiten.KeyDigit1) {
		g.bird.color = "red"
	}
	if ebiten.IsKeyPressed(ebiten.KeyDigit2) {
		g.bird.color = "yellow"
	}
	if ebiten.IsKeyPressed(ebiten.KeyDigit3) {
		g.bird.color = "green"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		switch g.bird.color {
		case "red":
			g.bird.color = "yellow"
			return
		case "yellow":
			g.bird.color = "green"
			return
		case "green":
			g.bird.color = "red"
			return
		}
	}
}

// 全てのボールを動かす
func (g *Game) ballsMove() {
	for i := range g.balls.ball {
		g.ballMove(i)
	}
}

// 1つのボールを動かす
func (g *Game) ballMove(index int) {
	g.balls.ball[index].y += g.balls.ball[index].speed
	g.balls.ball[index].speed += g.balls.ball[index].acceleration
}

// 鳥の翼を上げ下げする
func (g *Game) birdAnimation() {
	// cntの1桁目を取り出す
	cnt_single_digit := g.cnt % 10
	if cnt_single_digit == 0 {
		g.bird.state = "up"
	} else if cnt_single_digit == 5 {
		g.bird.state = "down"
	}
}

// 鳥とボール当たったか判定する
func (g *Game) hitBall(ball ball) bool {
	// 鳥とボールの距離dを求める
	birdmx := g.bird.x + g.bird.size/2
	birdmy := g.bird.y + g.bird.size/2
	ballmx := ball.x + ball.size/2
	ballmy := ball.y + ball.size/2
	dx := birdmx - ballmx
	dy := birdmy - ballmy
	dxE2 := math.Pow(dx, 2)
	dyE2 := math.Pow(dy, 2)
	d := math.Sqrt(dxE2 + dyE2)
	if d <= 60 { // 当たったらtrueを返す
		return true
	} else { // それ以外はfalseを返す
		return false
	}
}

// 鳥とボールが同じ色か判定する
func (g *Game) sameColor(ball ball) bool {
	if g.bird.color == ball.color {
		return true
	}
	return false
}

// 鳥とボールが当たったときの処理をする
func (g *Game) processHitBall() {
	for i, e := range g.balls.ball {
		if g.hitBall(e) && g.sameColor(e) { // ボールが当たって、かつ鳥が同じ色なら
			g.ballInit(i)
			g.hitScore++
			g.bird.speed += 0.01
		} else if g.hitBall(e) && !g.sameColor(e) { // ボールが当たって、かつ鳥が違う色なら
			g.bird.y += 5
		}
	}
}

// 鳥がスクリーンの中にいるか判定する
func (g *Game) isInScreen_Bird() bool {
	if g.bird.y > float64(screenHeight)-g.bird.size {
		return false
	}
	return true
}

// 鳥がスクリーンの外にいるときの処理をする
func (g *Game) processOutScreen_Bird() {
	if !g.isInScreen_Bird() {
		g.bird.state = "sleep"
		g.highscore = max(g.score, g.highscore)
		g.maxHitScore = max(g.hitScore, g.maxHitScore)
		g.mode = gameover
	}
}

// ボールがスクリーンの中にいるか判定する
func (g *Game) isInScreen_Ball(ball ball) bool {
	if ball.y > float64(screenWidth) {
		return false
	}
	return true
}

// ボールがスクリーンの外にあるときの処理をする
func (g *Game) processOutScreen_Balls() {
	for i, e := range g.balls.ball {
		if !g.isInScreen_Ball(e) {
			g.ballInit(i)
		}
	}
}

// 一定の時間が経つとボールの個数を追加する
func (g *Game) addBall() {
	if g.cnt%1250 == 1249 {
		s := ball{}
		g.balls.ball = append(g.balls.ball, s)
		g.ballInit(len(g.balls.ball) - 1)
		g.balls.number = len(g.balls.ball)
	}
}

// 一定の時間が経つとボールの加速度を増やす
func (g *Game) increaseBallAcceleration() {
	if g.cnt%1250 == 1249 {
		g.balls.accelerations = append(g.balls.accelerations, g.balls.accelerations[len(g.balls.accelerations)-1]+0.025)
	}
}
