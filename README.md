# QRAscii

A library to convert QRCode image to ASCII-art, it is useful when you want to 
display a QRCode in terminal.

Things this lib can NOT do:

* Generate QRCode from raw data.
* Get raw data from QRCode.

## Usage

```go

import (
    "github.com/deadblue/qrascii"
    "log"
    "os"

    // You need load image decoders by yourself
    _ "image/png"
)

func main() {
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
```

Then you will see the ASCII-art of the QRCode, like this:

```
█████████████████████████████████
██ ▄▄▄▄▄ ██▀  ███▄▄█▄▄ █ ▄▄▄▄▄ ██
██ █   █ █ █▀▀▀█▄█▄▀▄███ █   █ ██
██ █▄▄▄█ █ █ ▀▀   ▀▀ ▄ █ █▄▄▄█ ██
██▄▄▄▄▄▄▄█ █▄▀▄▀▄▀▄█ █ █▄▄▄▄▄▄▄██
██ █    ▄▀▀▀▀▄██▄█▀ ▀▀▀█   ▄▄█▀██
██▀▄ ▀▀▀▄▀██▄▄ █▀ ▄█▀█▀███▄▀█▀███
███ ▀ ▄▄▄██▀▀▄▄▀▀█▀▄█▀▀▀▀▀▀▄▄█▀██
██▀█ █▄ ▄██▀▀█  ▄▄▄██▀  ▀▀ ▄▄▀███
██▀ ▄ ▄▀▄▀█▄▀▀ █▄▄▀▄▀ █ ▀ ▀▄ █▀██
██ █ █▄▀▄▄▄  ▀▀█▀▄▄███ ▀ ▄▄█▄▀███
██▄█████▄█ ▀▄▄█▀▀▄   ▄ ▄▄▄ ▀   ██
██ ▄▄▄▄▄ █▀███  ▄ ▄▄█  █▄█ ▄▄█▀██
██ █   █ █ █▀▀▄▀▄▄▀▄█▀ ▄▄▄▄▀  ▀██
██ █▄▄▄█ █▄▄▀██▀█▄██▄▄  ▄ ▄ ▄ ███
██▄▄▄▄▄▄▄█▄▄██▄█▄▄▄▄█▄██▄▄▄█▄████
█████████████████████████████████
```

## License

MIT