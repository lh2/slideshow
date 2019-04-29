package main // import "entf.net/slideshow"

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strings"

	_ "github.com/hullerob/go.farbfeld"

	"github.com/ktye/ui"
	"github.com/ktye/ui/dpy"
	"github.com/nfnt/resize"
)

var win *ui.Window

func parseColor(bgColor string) (color.Color, error) {
	bgColor = strings.TrimPrefix(bgColor, "#")
	bgColor = strings.TrimPrefix(bgColor, "0x")
	b, err := hex.DecodeString(bgColor)
	if err != nil {
		return color.Black, err
	}

	return color.NRGBA{
		b[0],
		b[1],
		b[2],
		uint8(255),
	}, nil
}

func logerr(errch <-chan error) {
	for err := range errch {
		fmt.Fprintf(os.Stderr, "draw: %v\n", err)
	}
}

func readStdin(ch chan<- string) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		ch <- s.Text()
	}
	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "stdin: error reading: %v\n", err)
		os.Exit(1)
	}
}

func loadImage(p string) (image.Image, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	rgba := image.NewRGBA(image.Rectangle{image.ZP, b.Size()})
	draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)

	return rgba, nil
}

type imageDisplay struct {
	BGColor color.Color
	srcimg  image.Image
	img     image.Image
	newimg  bool
	lastsz  image.Rectangle
}

func (id *imageDisplay) SetImage(img image.Image) {
	id.srcimg = img
	id.newimg = true
}

func (id *imageDisplay) Draw(dst *image.RGBA, force bool) {
	if !force {
		return
	}
	r := dst.Bounds()
	draw.Draw(dst, r, image.NewUniform(id.BGColor), image.ZP, draw.Src)
	if id.srcimg != nil && (id.lastsz == image.ZR || id.lastsz != r || id.newimg) {
		id.img = resize.Thumbnail(uint(r.Dx()), uint(r.Dy()), id.srcimg, resize.Bicubic).(*image.RGBA)
		id.lastsz = r
		id.newimg = false
	}
	if id.img != nil {
		ir := id.img.Bounds()
		x := (r.Dx() - ir.Dx()) / 2
		y := (r.Dy() - ir.Dy()) / 2
		draw.Draw(dst, image.Rect(x, y, x+ir.Dx(), y+ir.Dy()), id.img, image.ZP, draw.Src)
	}
}
func (id *imageDisplay) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	return 0
}
func (id *imageDisplay) Key(r rune, code uint32, dir int, mod uint32) int {
	if code == 41 { // ESC
		go func() {
			win.Quit <- true
		}()
	}
	return 0
}

func main() {
	var bgColor string
	flag.StringVar(&bgColor, "bg", "#000000", "Background color as 6 digit hex value")
	flag.Parse()

	bgc, err := parseColor(bgColor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bg: invalid format\n")
		os.Exit(1)
	}

	win = ui.New(dpy.New(nil))
	imgdisp := &imageDisplay{}
	imgdisp.BGColor = bgc
	win.Top = imgdisp
	done := win.Run()

	in := make(chan string)
	go readStdin(in)

mainloop:
	for {
		select {
		case <-done:
			break mainloop
		case line := <-in:
			i, err := loadImage(line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "draw: unable to read image: %v\n", err)
				os.Exit(1)
			}
			imgdisp.SetImage(i)
			win.Draw(true)
		}
	}
}
