package main

import (
	"bytes"
	"fmt"
	"image"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type color int

const (
	Black color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func (c color) String() string {
	switch c {
	case Black:
		return "Black"
	case Red:
		return "Red"
	case Green:
		return "Green"
	case Yellow:
		return "Yellow"
	case Blue:
		return "Blue"
	case Magenta:
		return "Magenta"
	case Cyan:
		return "Cyan"
	case White:
		return "White"
	default:
		panic(fmt.Errorf("Unknown color: %v", int(c)))
	}

}

type RGB struct {
	R, G, B uint32
}

const MX = 0xffff

func (c color) Rgb() RGB {
	switch c {
	case Black:
		return RGB{0, 0, 0}
	case Red:
		return RGB{MX, 0, 0}
	case Green:
		return RGB{0, MX, 0}
	case Yellow:
		return RGB{MX, MX, 0}
	case Blue:
		return RGB{0, 0, MX}
	case Magenta:
		return RGB{MX, 0, MX}
	case Cyan:
		return RGB{0, MX, MX}
	case White:
		return RGB{MX, MX, MX}
	default:
		panic(fmt.Errorf("Unknown color: %v", int(c)))
	}
}

func (c color) Ascii() string {
	switch c {
	case Black:
		return " "
	case Red:
		return "\u001b[31m#\u001b[0m"
	case Green:
		return "\u001b[32m#\u001b[0m"
	case Yellow:
		return "\u001b[33m#\u001b[0m"
	case Blue:
		return "\u001b[34m#\u001b[0m"
	case Magenta:
		return "\u001b[35m#\u001b[0m"
	case Cyan:
		return "\u001b[36m#\u001b[0m"
	case White:
		return "\u001b[37m#\u001b[0m"
	default:
		panic(fmt.Errorf("Unknown color: %v", int(c)))
	}
}

func NearestColor(rgb RGB) color {
	col := Black
	cur := dist(rgb, col.Rgb())
	for _, c := range []color{Red, Green, Yellow, Blue, Magenta, Cyan, White} {
		if d := dist(rgb, c.Rgb()); d < cur {
			col = c
			cur = d
		}
	}
	// log.Printf("%v -> %v", rgb, col.String())
	return col
}

func dist(a, b RGB) uint32 {
	return absDiff(a.R, b.R) + absDiff(a.G, b.G) + absDiff(a.B, b.B)
}

func absDiff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}

const (
	// H = 32
	W = 120
)

func image2ascii(data []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	bounds := img.Bounds()
	H := int(float64(W) * float64(bounds.Max.Y-bounds.Min.Y) / float64(bounds.Max.X-bounds.Min.X))
	out := new(strings.Builder)
	for y := 0; y < H; y++ {
		py := bounds.Min.Y + int((float64(y)/float64(H))*float64(bounds.Max.Y-bounds.Min.Y))
		for x := 0; x < W; x++ {
			px := bounds.Min.X + int((float64(x)/float64(W))*float64(bounds.Max.X-bounds.Min.X))
			// log.Printf("px, py = %v, %v", px, py)
			c := img.At(px, py)
			r, g, b, _ := c.RGBA()
			a := NearestColor(RGB{r, g, b})
			fmt.Fprint(out, a.Ascii())
		}
		fmt.Fprintln(out)
	}
	return out.String(), nil
}
