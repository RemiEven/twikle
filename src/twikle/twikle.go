package main

import (
        "fmt"
		"image"
        "image/png"
		"image/draw"
		"image/color"
        "log"
        "os"
		"github.com/containous/flaeg"
		"github.com/nfnt/resize"
)

type Configuration struct {
	Image        string  `short:"i" description:"Image to transform"`
	Pattern      string  `short:"p" description:"Pattern file"`
	PatternScale uint `short:"s" description:"Scale of the pattern"`
}

//newDefaultConfiguration returns a pointer on Configuration with default values
func newDefaultPointersConfiguration() *Configuration {
    return newConfiguration()
}

//newConfiguration returns a pointer on an initialized configuration
func newConfiguration() *Configuration {
    return &Configuration{
        Input: "image.png",
        Pattern: "brick.png",
        PatternScale: 10,
    }
}

func main() {
    config := newConfiguration()

    rootCmd := &flaeg.Command{
    	Name: "twikle",
    	Description: `twikle is a program made to apply a pattern to an image`,
    	Config: config,
        DefaultPointersConfig: newDefaultPointersConfiguration(),
        Run: func() error {
            img := loadImage("/images/input/" + config.Image)

            result := image.NewRGBA(img.Bounds())
            draw.Draw(result, img.Bounds(), img, image.ZP, draw.Src)
            result = tile(result, config)

            dest, _ := os.Create("/images/output/" + config.Image)
            defer dest.Close()

            png.Encode(dest, result)
            return nil
        },
    }

    flaeg := flaeg.New(rootCmd, os.Args[1:])

    if err := flaeg.Run(); err != nil {
        fmt.Errorf("Error %s", err.Error())
    }
}

func tile(result *image.RGBA, config *Configuration) *image.RGBA {
	brickImage := loadImage("/images/pattern/" + config.Pattern)
	brickImage = resize.Resize(config.PatternScale, 0, brickImage, resize.MitchellNetravali)
	tileWidth := brickImage.Bounds().Max.X
	tileHeight := brickImage.Bounds().Max.Y

	tileArea := tileWidth * tileHeight

	lineTileWidth := result.Bounds().Max.X / tileWidth
	columnTileHeight := result.Bounds().Max.Y / tileHeight

	extraWidth := result.Bounds().Max.X % tileWidth
	extraHeight := result.Bounds().Max.Y % tileHeight

	croppedRectangle := result.Bounds()
	croppedRectangle.Min.X += extraWidth / 2
	croppedRectangle.Min.Y += extraHeight / 2
	croppedRectangle.Max.X -= extraWidth / 2
	croppedRectangle.Max.Y -= extraHeight / 2
	croppedImage := result.SubImage(croppedRectangle).(*image.RGBA)

	r := brickImage.Bounds()
	for y := 0; y < columnTileHeight; y++ {
		for x := 0; x < lineTileWidth; x++ {
			r.Min.X = extraWidth / 2 + x * tileWidth
			r.Max.X = r.Min.X + tileWidth

			r.Min.Y = extraHeight / 2 + y * tileHeight
			r.Max.Y = r.Min.Y + tileHeight

			var averageR uint64 = 0
			var averageG uint64 = 0
			var averageB uint64 = 0
			var averageA uint64 = 0

            // Why not just use a resize ?
			for py := 0; py < tileHeight; py++ {
				for px := 0; px < tileWidth; px++ {
					r, g, b, a := croppedImage.At(r.Min.X + px, r.Min.Y + py).RGBA()
					averageR += uint64(r)
					averageG += uint64(g)
					averageB += uint64(b)
					averageA += uint64(a)
				}
			}
			averageR /= uint64(tileArea) * 256
			averageG /= uint64(tileArea) * 256
			averageB /= uint64(tileArea) * 256
			averageA /= uint64(tileArea) * 256

			averageColor := color.RGBA{uint8(averageR), uint8(averageG), uint8(averageB), uint8(averageA)}
			draw.Draw(croppedImage, r, &image.Uniform{averageColor}, image.ZP, draw.Src)

			tileAlpha := color.Alpha{uint8(averageA)}
			draw.DrawMask(croppedImage, r, brickImage, image.ZP, &image.Uniform{tileAlpha}, image.ZP, draw.Over)
		}
	}
	return croppedImage
}

func loadImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
