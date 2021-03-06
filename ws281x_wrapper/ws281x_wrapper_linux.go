// +build linux,cgo darwin,cgo

package ws281x_wrapper

import (
	"errors"
	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

// Engine represents a wrapper around the ws281x library.
// Holds the state of the leds, the leds count and a reference to the actual
// ws281x library instance
type Engine struct {
	engine   *ws281x.WS2811
	leds     []uint32
	ledCount int
}

// Init initializes a new instance of the ws281x library
func Init(pin int, ledCount int, brightness int) (*Engine, error) {
	// Initialize ws281x engine
	opt := ws281x.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCount
	opt.Channels[0].GpioPin = pin
	ws, err := ws281x.MakeWS2811(&opt)
	if err != nil {
		return nil, err
	}
	err = ws.Init()
	if err != nil {
		return nil, err
	}
	engine := &Engine{
		leds:     make([]uint32, ledCount),
		engine:   ws,
		ledCount: ledCount,
	}
	return engine, nil
}

// Fini does cleanup operations
func (ws *Engine) Fini() {
	ws.engine.Fini()
}

// Clear resets all the leds (turns them off by setting their color to black)
func (ws *Engine) Clear() error {
	ws.leds = make([]uint32, ws.ledCount)
	return ws.Render()
}

// Render renders the colors saved on the leds array onto the led strip
func (ws *Engine) Render() error {
	ws.engine.SetLedsSync(0, ws.leds)
	return ws.engine.Render()
}

// SetLedColor changes the color of the led in the specified index
func (ws *Engine) SetLedColor(index int, r uint8, g uint8, b uint8) error {
	if index >= len(ws.leds) || index < 0 {
		return errors.New("Invalid led index")
	}
	// WRGB
	color := uint32(0xff)<<24 | uint32(r)<<16 | uint32(g)<<8 | uint32(b)
	ws.leds[index] = color
	return nil
}
