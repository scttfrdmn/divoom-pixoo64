package pixoo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"os/exec"
)

// CurlClient uses curl command to bypass macOS security restrictions
type CurlClient struct {
	host string
}

func NewCurlClient(host string) *CurlClient {
	return &CurlClient{host: host}
}

func (c *CurlClient) post(command interface{}) error {
	url := fmt.Sprintf("http://%s/post", c.host)

	data, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("marshal command: %w", err)
	}

	// Use curl command which is already trusted by macOS
	cmd := exec.Command("curl", "-s", "-X", "POST", url,
		"-H", "Content-Type: application/json",
		"-d", string(data),
		"-m", "5")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("curl command: %w (output: %s)", err, string(output))
	}

	return nil
}

func (c *CurlClient) SetBrightness(brightness int) error {
	if brightness < 0 || brightness > 100 {
		return fmt.Errorf("brightness must be between 0 and 100")
	}

	command := map[string]interface{}{
		"Command":    "Channel/SetBrightness",
		"Brightness": brightness,
	}

	return c.post(command)
}

func (c *CurlClient) ClearScreen() error {
	command := map[string]interface{}{
		"Command": "Draw/ResetHttpGifId",
	}
	return c.post(command)
}

func (c *CurlClient) DrawImage(img image.Image) error {
	bounds := img.Bounds()
	if bounds.Dx() != 64 || bounds.Dy() != 64 {
		return fmt.Errorf("image must be 64x64 pixels")
	}

	// Convert image to RGB data
	pixels := make([]int, 64*64)
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels[y*64+x] = int((r>>8)<<16 | (g>>8)<<8 | (b >> 8))
		}
	}

	command := map[string]interface{}{
		"Command":   "Draw/SendHttpGif",
		"PicNum":    1,
		"PicWidth":  64,
		"PicOffset": 0,
		"PicID":     0,
		"PicSpeed":  1000,
		"PicData":   pixels,
	}

	// For large payloads, write to temp file
	data, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("marshal command: %w", err)
	}

	url := fmt.Sprintf("http://%s/post", c.host)
	cmd := exec.Command("curl", "-s", "-X", "POST", url,
		"-H", "Content-Type: application/json",
		"--data-binary", "@-",
		"-m", "10")

	cmd.Stdin = bytes.NewReader(data)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("curl command: %w (output: %s)", err, string(output))
	}

	return nil
}

func (c *CurlClient) DrawText(text string, r, g, b uint8) error {
	command := map[string]interface{}{
		"Command":    "Draw/SendHttpText",
		"TextId":     1,
		"x":          0,
		"y":          0,
		"dir":        0,
		"font":       2,
		"TextWidth":  64,
		"speed":      50,
		"TextString": text,
		"color":      fmt.Sprintf("#%02x%02x%02x", r, g, b),
		"align":      1,
	}

	return c.post(command)
}
