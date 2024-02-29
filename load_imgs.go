package main

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed imgs/*.png
var images embed.FS

const (
	birdSize float64 = 85
	ballSize float64 = 57
)

var (
	backgroundImage *ebiten.Image
	birdImages      map[bird_Img_Structure]*ebiten.Image
	ballImages      map[string]*ebiten.Image
	colors          = []string{"red", "yellow", "green"}
	states          = []string{"up", "down", "sleep"}
)

type bird_Img_Structure struct {
	color, state string
}

func init() {
	fname := "imgs/background.png"
	file, err := images.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	backgroundImage = ebiten.NewImageFromImage(img)

	birdImages = make(map[bird_Img_Structure]*ebiten.Image)

	for _, e := range colors[0:3] {
		for _, ee := range states {
			fname = fmt.Sprintf("imgs/bird_%s_%s.png", e, ee)
			file, err = images.Open(fname)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			img, _, err = image.Decode(file)
			if err != nil {
				log.Fatal(err)
			}
			birdImages[bird_Img_Structure{e, ee}] = ebiten.NewImageFromImage(img)
		}
	}

	ballImages = make(map[string]*ebiten.Image)
	for _, e := range colors {
		fname := fmt.Sprintf("imgs/ball_%s.png", e)
		file, err := images.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		ballImages[e] = ebiten.NewImageFromImage(img)
	}
}
