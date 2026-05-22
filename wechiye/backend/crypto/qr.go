package crypto

import (
	"bytes"
	"image/png"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(data string) ([]byte, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, qr.Image(256))
	return buf.Bytes(), err
}