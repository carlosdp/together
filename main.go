package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
)

var strokeFont *truetype.Font
var posterFont *truetype.Font

func init() {
	fontBytes, err := ioutil.ReadFile("brush_strokes.ttf")
	if err != nil {
		panic(err)
	}

	strokeFont, err = truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	spb, err := ioutil.ReadFile("SansPosterBold.ttf")
	if err != nil {
		panic(err)
	}

	posterFont, err = truetype.Parse(spb)
	if err != nil {
		panic(err)
	}
}

func main() {
	rightFile, err := os.Open("p1.jpg")
	if err != nil {
		panic(err)
	}

	rightPortrait, err := jpeg.Decode(rightFile)
	if err != nil {
		panic(err)
	}

	leftFile, err := os.Open("p2.jpg")
	if err != nil {
		panic(err)
	}

	leftPortrait, err := jpeg.Decode(leftFile)
	if err != nil {
		panic(err)
	}

	rightImage := resize.Resize(400, 600, rightPortrait, resize.Bilinear)
	leftImage := resize.Resize(400, 600, leftPortrait, resize.Bilinear)

	t, err := os.Open("template.png")
	if err != nil {
		panic(err)
	}
	mask, err := png.Decode(t)
	if err != nil {
		panic(err)
	}

	newImage := image.NewRGBA(rightImage.Bounds())
	draw.Draw(newImage, rightImage.Bounds(), rightImage, image.ZP, draw.Src)
	draw.DrawMask(newImage, leftImage.Bounds(), leftImage, image.ZP, mask, image.ZP, draw.Over)

	ctx := freetype.NewContext()
	ctx.SetFont(strokeFont)
	ctx.SetFontSize(84)
	ctx.SetSrc(image.Black)
	ctx.SetDst(newImage)
	ctx.SetClip(newImage.Bounds())
	_, err = ctx.DrawString("TOGETHER", freetype.Pt(40, 200))
	if err != nil {
		panic(err)
	}

	ctx.SetFont(posterFont)
	ctx.SetFontSize(12)
	ctx.SetSrc(image.White)
	ctx.SetDst(newImage)
	ctx.SetClip(newImage.Bounds())
	_, err = ctx.DrawString("#VOTETOGETHER", freetype.Pt(150, 580))
	if err != nil {
		panic(err)
	}

	f, err := os.Create("newimage.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = png.Encode(f, newImage)
	if err != nil {
		panic(err)
	}
}
