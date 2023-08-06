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

// ToAscii converts dot map to ascii art, the result is an unicode text.
func (q *QRCode) ToAscii(margin int) string {
	buf := make([]rune, 0)
	for y := 0 - margin; y < q.size+margin; y += 2 {
		for x := 0 - margin; x < q.size+margin; x++ {
			u, d := 0, 0
			if q.get(x, y) {
				u = 1
			}
			if q.get(x, y+1) {
				d = 1
			}
			buf = append(buf, blockChars[d*2+u])
		}
		buf = append(buf, '\n')
	}
	return string(buf)
}
