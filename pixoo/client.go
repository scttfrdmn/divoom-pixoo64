package pixoo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"time"
)

const (
	DefaultPort = 80
)

type Client struct {
	host       string
	httpClient *http.Client
}

func NewClient(host string) *Client {
	// Use default transport - custom transports can cause permission issues on macOS
	return &Client{
		host: host,
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // Increased timeout for image uploads
		},
	}
}

func (c *Client) post(command interface{}) error {
	url := fmt.Sprintf("http://%s/post", c.host)

	data, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("marshal command: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// SetBrightness sets the screen brightness (0-100)
func (c *Client) SetBrightness(brightness int) error {
	if brightness < 0 || brightness > 100 {
		return fmt.Errorf("brightness must be between 0 and 100")
	}

	command := map[string]interface{}{
		"Command": "Channel/SetBrightness",
		"Brightness": brightness,
	}

	return c.post(command)
}

// SetChannel switches to a specific channel
// 0 = Faces/Clock, 1 = Cloud Channel, 2 = Visualizer, 3 = Custom
func (c *Client) SetChannel(channel int) error {
	command := map[string]interface{}{
		"Command":     "Channel/SetIndex",
		"SelectIndex": channel,
	}
	return c.post(command)
}

// ClearScreen clears the display
func (c *Client) ClearScreen() error {
	command := map[string]interface{}{
		"Command": "Draw/ResetHttpGifId",
	}
	return c.post(command)
}

// DrawImage sends a 64x64 image to the display
// The image is converted to base64-encoded RGB format expected by Pixoo
func (c *Client) DrawImage(img image.Image) error {
	bounds := img.Bounds()
	if bounds.Dx() != 64 || bounds.Dy() != 64 {
		return fmt.Errorf("image must be 64x64 pixels")
	}

	// Reset GIF state first
	resetCmd := map[string]interface{}{
		"Command": "Draw/ResetHttpGifId",
	}
	if err := c.post(resetCmd); err != nil {
		return fmt.Errorf("reset gif: %w", err)
	}

	// Convert image to RGB bytes (R,G,B,R,G,B,...)
	// 64x64 pixels * 3 bytes = 12,288 bytes
	pixelBytes := make([]byte, 64*64*3)
	idx := 0
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixelBytes[idx] = byte(r >> 8)     // R
			pixelBytes[idx+1] = byte(g >> 8)   // G
			pixelBytes[idx+2] = byte(b >> 8)   // B
			idx += 3
		}
	}

	// Base64 encode the pixel data
	encodedData := base64.StdEncoding.EncodeToString(pixelBytes)

	command := map[string]interface{}{
		"Command":   "Draw/SendHttpGif",
		"PicID":     1,
		"PicNum":    1,
		"PicOffset": 0,
		"PicWidth":  64,
		"PicSpeed":  1000,
		"PicData":   encodedData,
	}

	return c.post(command)
}

// DrawText displays text on the screen
func (c *Client) DrawText(text string, r, g, b uint8) error {
	command := map[string]interface{}{
		"Command":    "Draw/SendHttpText",
		"TextId":     1,
		"x":          0,
		"y":          24,  // Center vertically
		"dir":        0,   // 0 = left scroll
		"font":       2,
		"TextWidth":  64,
		"speed":      0,   // 0 = static, >0 = scrolling
		"TextString": text,
		"color":      fmt.Sprintf("#%02x%02x%02x", r, g, b),
		"align":      2,  // 2 = center
	}

	return c.post(command)
}

// CreateImage creates a blank 64x64 image
func CreateImage() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, 64, 64))
}

// FillRect fills a rectangle on the image
func FillRect(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	for y := y1; y < y2 && y < 64; y++ {
		for x := x1; x < x2 && x < 64; x++ {
			img.Set(x, y, c)
		}
	}
}

// DrawTextOnImage draws simple text on image (basic 5x7 font)
func DrawTextOnImage(img *image.RGBA, text string, x, y int, c color.Color) {
	// For now, just a placeholder - you can implement a bitmap font or use a library
	// This is a simple implementation that draws text at given position
	for i, ch := range text {
		drawChar(img, ch, x+i*6, y, c)
	}
}

func drawChar(img *image.RGBA, ch rune, x, y int, c color.Color) {
	// Simple 5x7 pixel font for basic characters (0-9, A-Z, symbols)
	// This is a basic implementation - you can expand this
	patterns := getCharPattern(ch)

	for row := 0; row < len(patterns); row++ {
		for col := 0; col < 5; col++ {
			if patterns[row]&(1<<(4-col)) != 0 {
				if x+col < 64 && y+row < 64 && x+col >= 0 && y+row >= 0 {
					img.Set(x+col, y+row, c)
				}
			}
		}
	}
}

func getCharPattern(ch rune) []byte {
	// 5x7 bitmap font patterns
	switch ch {
	case '0':
		return []byte{0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E}
	case '1':
		return []byte{0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E}
	case '2':
		return []byte{0x0E, 0x11, 0x01, 0x02, 0x04, 0x08, 0x1F}
	case '3':
		return []byte{0x1F, 0x02, 0x04, 0x02, 0x01, 0x11, 0x0E}
	case '4':
		return []byte{0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02}
	case '5':
		return []byte{0x1F, 0x10, 0x1E, 0x01, 0x01, 0x11, 0x0E}
	case '6':
		return []byte{0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E}
	case '7':
		return []byte{0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08}
	case '8':
		return []byte{0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E}
	case '9':
		return []byte{0x0E, 0x11, 0x11, 0x0F, 0x01, 0x02, 0x0C}
	case '%':
		return []byte{0x18, 0x19, 0x02, 0x04, 0x08, 0x13, 0x03}
	case ':':
		return []byte{0x00, 0x00, 0x0C, 0x0C, 0x00, 0x0C, 0x0C}
	case ' ':
		return []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 'C':
		return []byte{0x0E, 0x11, 0x10, 0x10, 0x10, 0x11, 0x0E}
	case 'P':
		return []byte{0x1E, 0x11, 0x11, 0x1E, 0x10, 0x10, 0x10}
	case 'U':
		return []byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E}
	case 'M':
		return []byte{0x11, 0x1B, 0x15, 0x15, 0x11, 0x11, 0x11}
	case 'E':
		return []byte{0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x1F}
	case 'G':
		return []byte{0x0E, 0x11, 0x10, 0x17, 0x11, 0x11, 0x0F}
	case 'B':
		return []byte{0x1E, 0x11, 0x11, 0x1E, 0x11, 0x11, 0x1E}
	default:
		return []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	}
}
