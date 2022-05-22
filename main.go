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
	"strings"
	"time"
)

const (
	fontSize    = 24
	fontPath    = "font/SourceHanSansCN-Light.otf"
	outputPath  = "out/out.png"
	lineSpacing = 0
)

func measureString() (int, int) {

	return 0, 0
}

type textLine struct {
	text string
}

type TextDrawer struct {
	lines []textLine
}

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

func addLabel(img *image.RGBA, x int, label string) {
	col := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: loadFontFace(),
	}
	runes := []rune(label)
	last := 0
	row := 1
	for i := range runes {
		if d.MeasureString(string(runes[last:i])).Round() > 800 {
			d.Dot = fixed.Point26_6{X: fixed.I(x), Y: fixed.I(fontSize*row + (row-1)*lineSpacing)}
			d.DrawString(string(runes[last : i-1]))
			last = i - 1
			row += 1
		}
	}
	if left := runes[last:]; len(left) > 0 {
		d.Dot = fixed.Point26_6{X: fixed.I(x), Y: fixed.I(fontSize*row + (row-1)*lineSpacing)}
		d.DrawString(string(left))
	}
}

func main() {
	start := time.Now()
	img := image.NewRGBA(image.Rect(0, 0, 800, 1000))
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
	ss, err := ioutil.ReadFile("text.txt")
	if err != nil {
		panic(err)
	}
	addLabel(img, 0, strings.NewReplacer(
		"\n", "",
		"\t", "",
		"\r", "",
	).Replace(string(ss)))

	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		panic(err)
	}
	log.Println("cost", time.Now().UnixMilli()-start.UnixMilli(), "ms")
}
