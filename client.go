package pixgo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// NewClient instantiates a new client.  address is the IP address
// of the Pixoo device, and size is the size of the display (e.g. 16, 32, 64).
func NewClient(address string, size int) *Client {
	return &Client{
		Address:    address,
		Size:       size,
		buffer:     make([]int, size*size*3),
		reqCounter: 1,
	}
}

type Client struct {
	Address    string
	Size       int
	buffer     []int
	reqCounter int
}

func (c *Client) doRequest(body interface{}) error {
	marshalled, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	_, err = http.Post("http://"+c.Address+"/post", "application/json", bytes.NewBuffer(marshalled))
	if err != nil {
		return fmt.Errorf("error making request to Pixoo: %w", err)
	}

	return nil
}

// PixelCount returns the number of pixels in the display.
func (c *Client) PixelCount() int {
	return c.Size * c.Size
}

// SetBrightness changes the brightness of the display.
func (c *Client) SetBrightness(brightness int) error {
	reqBody := map[string]interface{}{
		"Command":    "Channel/SetBrightness",
		"Brightness": brightness,
	}

	return c.doRequest(reqBody)
}

// SetPowerStatus turns the device on or off; true for on, false for off.
func (c *Client) SetPowerStatus(status bool) error {
	var statusInt = 0
	if status {
		statusInt = 1
	}

	reqBody := map[string]interface{}{
		"Command": "Channel/OnOffScreen",
		"OnOff":   statusInt,
	}

	return c.doRequest(reqBody)
}

// TurnOn is shorthand for SetPowerStatus(true).
func (c *Client) TurnOn() error {
	return c.SetPowerStatus(true)
}

// TurnOff is shorthand for SetPowerStatus(false).
func (c *Client) TurnOff() error {
	return c.SetPowerStatus(false)
}

// GetIndex returns the starting index of the pixel in the buffer.
func (c *Client) GetIndex(x int, y int) int {
	return (x + (y * c.Size)) * 3
}

// DrawCharacter draws a character at the given position.
func (c *Client) DrawCharacter(char rune, x, y, red, green, blue int) error {
	charData, ok := chars[char]
	if !ok {
		return errors.New("unrecognized character")
	}

	for i, v := range charData {
		if v == 1 {
			localX := i % 3
			localY := int(i / 3)
			c.SetPixel(x+localX, y+localY, red, green, blue)
		}
	}
	return nil
}

// Fill fills the entire display with the given color.
func (c *Client) Fill(red, green, blue int) {
	for x := 0; x < c.Size; x++ {
		for y := 0; y < c.Size; y++ {
			c.SetPixel(x, y, red, green, blue)
		}
	}
}

// SetPixel sets the color of a pixel at a given coordinate.
func (c *Client) SetPixel(x, y, red, green, blue int) {
	index := c.GetIndex(x, y)
	c.buffer[index] = red
	c.buffer[index+1] = green
	c.buffer[index+2] = blue
}

// Reset resets the display to its default state, removing
// any items in the display.
func (c *Client) Reset() error {
	reqBody := map[string]interface{}{
		"Command": "Draw/ResetHttpGifId",
	}

	return c.doRequest(reqBody)
}

// Push flushes the buffer to the Pixoo device.
func (c *Client) Push() error {
	byteArray := make([]byte, len(c.buffer))
	for i, v := range c.buffer {
		byteArray[i] = byte(v)
	}

	encoded := base64.StdEncoding.EncodeToString(byteArray)
	reqBody := map[string]interface{}{
		"Command":   "Draw/SendHttpGif",
		"PicNum":    1,
		"PicWidth":  c.Size,
		"PicOffset": 0,
		"PicID":     c.reqCounter,
		"PicSpeed":  1000,
		"PicData":   encoded,
	}
	c.reqCounter += 1

	return c.doRequest(reqBody)
}
