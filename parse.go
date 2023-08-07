package qrascii

import (
	"image"
	"io"
)

type _Parser struct {
	im image.Image

	// The bound of QRCode
	bounds image.Rectangle
	// Size of a dot
	dotSize int
	//
	dotThreshold int

	// QRCode size
	Size int
}

func (p *_Parser) isBlack(x, y int) bool {
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

func (p *_Parser) locate() (err error) {
	p.bounds.Min.X, p.bounds.Min.Y = -1, -1
	p.bounds.Max.X, p.bounds.Max.Y = -1, -1
	// Locate QRcode
	ib := p.im.Bounds()
	// Lookup left-top corner
	for y := ib.Min.Y; y < ib.Max.Y; y++ {
		var found bool
		for x := ib.Min.X; x < ib.Max.X; x++ {
			found = p.isBlack(x, y)
			if found {
				p.bounds.Min.X, p.bounds.Min.Y = x, y
				break
			}
		}
		if found {
			break
		}
	}
	if p.bounds.Min.X < 0 {
		return ErrInvalidQrcode
	}
	// Lookup right-top corner
	for x := ib.Max.X - 1; x > p.bounds.Min.X; x-- {
		if p.isBlack(x, p.bounds.Min.Y) {
			p.bounds.Max.X = x + 1
			break
		}
	}
	if p.bounds.Max.X < 0 {
		return ErrInvalidQrcode
	}
	// Lookup left-bottom corner
	for y := ib.Max.Y - 1; y > p.bounds.Min.Y; y-- {
		if p.isBlack(p.bounds.Min.X, y) {
			p.bounds.Max.Y = y + 1
			break
		}
	}
	if p.bounds.Max.Y < 0 {
		return ErrInvalidQrcode
	}
	if p.bounds.Dx() != p.bounds.Dy() {
		return ErrInvalidQrcode
	}
	return
}

func (p *_Parser) measure() (err error) {
	for i := 0; i < p.bounds.Dx(); i++ {
		if !p.isBlack(p.bounds.Min.X+i, p.bounds.Min.Y+i) {
			p.dotSize = i
			break
		}
	}
	if p.dotSize == 0 {
		return ErrInvalidQrcode
	}
	if p.bounds.Dx()%p.dotSize != 0 {
		return ErrInvalidQrcode
	}
	p.Size = p.bounds.Dx() / p.dotSize
	if (p.Size-21)%4 != 0 {
		return ErrInvalidQrcode
	}
	p.dotThreshold = p.dotSize * p.dotSize * 8 / 10
	return
}

func (p *_Parser) Parse() (err error) {
	if err = p.locate(); err != nil {
		return
	}
	return p.measure()
}

func (p *_Parser) IsDotBlack(row, col int) bool {
	blackCount := 0
	for y := 0; y < p.dotSize; y++ {
		py := p.bounds.Min.Y + row*p.dotSize + y
		for x := 0; x < p.dotSize; x++ {
			px := p.bounds.Min.X + col*p.dotSize + x
			if p.isBlack(px, py) {
				blackCount += 1
			}
		}
	}
	return blackCount >= p.dotThreshold
}

// ParseImage parses QRCode image.
func ParseImage(im image.Image) (q *QRCode, err error) {
	p := &_Parser{im: im}
	if err = p.Parse(); err != nil {
		return
	}
	// Make QRcode
	q = &QRCode{
		size:   p.Size,
		matrix: make([]bool, p.Size*p.Size),
	}
	// Convert dot map to matrix
	for row := 0; row < q.size; row++ {
		for col := 0; col < q.size; col++ {
			q.matrix[row*q.size+col] = p.IsDotBlack(row, col)
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
