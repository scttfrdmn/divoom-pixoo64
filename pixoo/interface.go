package pixoo

import "image"

// PixooClient interface for different client implementations
type PixooClient interface {
	SetBrightness(brightness int) error
	ClearScreen() error
	DrawImage(img image.Image) error
	DrawText(text string, r, g, b uint8) error
}
