package qrascii

import (
	"image"
	_ "image/png"
	"io"
)

type parser struct {
	im image.Image
}

func (p *parser) isBlack(x, y int) bool {
	// Get color at point
	c := p.im.At(x, y)
	// Check color
	r, g, b, a := c.RGBA()
	if a == 0 {
		return false
	}
	gray := (19595*r + 38470*g + 7471*b + 1<<15) >> 24
	return gray < 127
}

func (p *parser) locateQrcode(bounds *image.Rectangle) (err error) {
	bounds.Min.X, bounds.Min.Y = -1, -1
	bounds.Max.X, bounds.Max.Y = -1, -1
	// Locate QRcode
	ib := p.im.Bounds()
	// Lookup left-top corner
	for y := ib.Min.Y; y < ib.Max.Y; y++ {
		var found bool
		for x := ib.Min.X; x < ib.Max.X; x++ {
			found = p.isBlack(x, y)
			if found {
				bounds.Min.X, bounds.Min.Y = x, y
				break
			}
		}
		if found {
			break
		}
	}
	if bounds.Min.X < 0 {
		return ErrInvalidQrcode
	}
	// Lookup right-top corner
	for x := ib.Max.X - 1; x > bounds.Min.X; x-- {
		if p.isBlack(x, bounds.Min.Y) {
			bounds.Max.X = x + 1
			break
		}
	}
	if bounds.Max.X < 0 {
		return ErrInvalidQrcode
	}
	// Lookup left-bottom corner
	for y := ib.Max.Y - 1; y > bounds.Min.Y; y-- {
		if p.isBlack(bounds.Min.X, y) {
			bounds.Max.Y = y + 1
			break
		}
	}
	if bounds.Max.Y < 0 {
		err = ErrInvalidQrcode
	}
	return
}

func (p *parser) measureDotSize(bounds *image.Rectangle) (int, error) {
	for i := 0; i < bounds.Dx(); i++ {
		if !p.isBlack(bounds.Min.X+i, bounds.Min.Y+i) {
			return i, nil
		}
	}
	return 0, ErrInvalidQrcode
}

func (p *parser) isDotBlack(o image.Point, size int, threshold int) bool {
	blackCount := 0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if p.isBlack(o.X+x, o.Y+y) {
				blackCount += 1
			}
		}
	}
	return blackCount >= threshold
}

// ParseImage parses QRCode image.
func ParseImage(im image.Image) (q *QRCode, err error) {
	p := &parser{im: im}
	// Locate qrcode
	qb := &image.Rectangle{}
	if err = p.locateQrcode(qb); err != nil {
		return
	}
	// Measure dot size
	var dotSize int
	if dotSize, err = p.measureDotSize(qb); err != nil {
		return
	}
	// Make QRcode
	q = &QRCode{
		size: qb.Dx() / dotSize,
	}
	q.matrix = make([]bool, q.size*q.size)
	// Convert dot map to matrix
	threshold := dotSize * dotSize * 8 / 10
	for i := 0; i < q.size; i++ {
		for j := 0; j < q.size; j++ {
			o := image.Point{
				X: qb.Min.X + j*dotSize,
				Y: qb.Min.Y + i*dotSize,
			}
			q.matrix[i*q.size+j] = p.isDotBlack(o, dotSize, threshold)
		}
	}
	return
}

// Parse parses QRCode image from an io.Reader.
func Parse(r io.Reader) (q *QRCode, err error) {
	// Decode image
	var im image.Image
	if im, _, err = image.Decode(r); err == nil {
		q, err = ParseImage(im)
	}
	return
}
