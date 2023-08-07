package qrascii

// QRCode holds the dot map of a QRCode.
type QRCode struct {
	size   int
	matrix []bool
}

func (q *QRCode) get(x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}
	if x >= q.size || y >= q.size {
		return false
	}
	return q.matrix[y*q.size+x]
}

// ToAscii converts dot map to ascii art with options.
func (q *QRCode) ToAscii(opts Options) string {
	if opts.Margin == 0 {
		opts.Margin = 2
	}
	buf := make([]rune, 0)
	for y := 0 - opts.Margin; y < q.size+opts.Margin; y += 2 {
		for x := 0 - opts.Margin; x < q.size+opts.Margin; x++ {
			u, d := 0, 0
			if q.get(x, y) {
				u = 1
			}
			if q.get(x, y+1) {
				d = 1
			}
			index := (d << 1) | u
			if opts.Invert {
				index = 3 ^ index
			}
			buf = append(buf, blockUnicode[index])
		}
		buf = append(buf, '\n')
	}
	return string(buf)
}

// Generate ascii-art with default options
func (q *QRCode) String() string {
	return q.ToAscii(Options{
		Margin: 2,
	})
}
