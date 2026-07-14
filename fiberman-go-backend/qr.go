package fiberman

import (
	"encoding/base64"
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

const defaultQRSize = 320

type QRCodeService struct{}

func NewQRCodeService() *QRCodeService {
	return &QRCodeService{}
}

func (s *QRCodeService) Generate(value string, requestedSize *int) (*InvoiceQRCodeResponse, error) {
	size := defaultQRSize
	if requestedSize != nil {
		size = *requestedSize
	}

	png, err := qrcode.Encode(value, qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice QR code: %w", err)
	}

	pngBase64 := base64.StdEncoding.EncodeToString(png)
	return &InvoiceQRCodeResponse{
		Value:     value,
		Size:      size,
		PNGBase64: pngBase64,
		DataURL:   "data:image/png;base64," + pngBase64,
	}, nil
}
