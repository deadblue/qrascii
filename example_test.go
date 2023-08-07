package qrascii

import (
	_ "image/png"
	"log"
	"os"
)

func ExampleParse() {
	qrFile := "/path/to/qrcode.png"
	f, err := os.Open(qrFile)
	if err != nil {
		log.Fatalf("Load QRCode image failed: %s", err)
	}
	defer f.Close()
	if qr, err := Parse(f); err != nil {
		log.Fatalf("Parse QRCode image failed: %s", err)
	} else {
		print(qr.String())
	}
}
