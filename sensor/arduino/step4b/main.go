package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/buzzer"
)

var (
	green  = machine.D12
	red    = machine.D10
	button = machine.D11
	touch  = machine.D9
	bzrPin = machine.D8
)

func main() {
	green.Configure(machine.PinConfig{Mode: machine.PinOutput})
	red.Configure(machine.PinConfig{Mode: machine.PinOutput})
	button.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	touch.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	bzrPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	bzr := buzzer.New(bzrPin)

	// morse code: tinygo
	delay := []int{
		300,
		0,
		100,
		100,
		0,
		300,
		100,
		0,
		300,
		100,
		300,
		300,
		0,
		300,
		300,
		100,
		0,
		300,
		300,
		300,
	}

	for _, d := range delay {
		green.High()
		red.High()
		bzr.On()
		time.Sleep(time.Millisecond * time.Duration(d))

		green.Low()
		red.Low()
		bzr.Off()
		time.Sleep(time.Millisecond * 100)

		if d == 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}

}
