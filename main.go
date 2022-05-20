package main

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

const (
	fontSize   = 48
	fontPath   = "font/SourceHanSansCN-Regular.otf"
	outputPath = "out/out.png"
)

func loadFontFace() font.Face {
	file, err := ioutil.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}
	parse, err := opentype.Parse(file)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(parse, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}
	return face
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: loadFontFace(),
		Dot:  point,
	}
	for _, i := range label {
		s := string(i)
		bounds, _ := d.BoundString(s)
		log.Println(s, "height:", bounds.Max.Y.Ceil()-bounds.Min.Y.Floor(), "width:", bounds.Max.X.Ceil()-bounds.Min.X.Floor())
	}
	d.DrawString(label)
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	for x := 0; x < img.Rect.Max.X; x++ {
		for y := 0; y < img.Rect.Max.Y; y++ {
			img.Set(x, y, color.RGBA{
				R: 255,
				G: 255,
				B: 255,
				A: 255,
			})
		}
	}
	addLabel(img, 0, fontSize, "Hello Go 那没事了")

	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		panic(err)
	}
}
