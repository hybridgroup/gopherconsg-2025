package main

import (
	"machine"
	"time"
)

// This Go program implements a simple 2-bit binary adder using Adruino microcontroller and LEDs.
// The user interacts with the system via buttons and touch inputs to select two 2-bit binary numbers,
// which are then added, with the result displayed on three LEDs.

// Key Features:

// Hardware Components:

// LEDs (Red, Green, Blue): Used to represent numbers (0-7) as binary.
// Buttons: For user input to cycle through states and confirm selections.
// Touch Sensor: Used to increment the binary input values.

// Program Flow:

// Initialization: All LEDs are turned on to indicate readiness.

// Input Capture:
// First 2-bit operand (input1) is selected using the touch sensor.
// After confirmation, the second operand (input2) is selected.
// Addition Display: The sum of input1 and input2 is calculated and displayed on the LEDs.
// Reset: Pressing the button resets the system to accept new inputs.

// User Interaction:

// Button presses and touch inputs guide the program through different states:
// Selecting the first operand.
// Selecting the second operand.
// Displaying the sum.
// Resetting the program for a new operation.

// To run: tinygo flash -target=nano-rp2040 main.go

var (
	green  = machine.D12
	red    = machine.D10
	button = machine.D11
	touch  = machine.D9
	blue   = machine.D8
)

// numberToLED: Maps a number (0-7) to a binary LED representation.
// red LED represents the least significant bit. blue LED represents
// the most significant bit.
func numberToLED(i uint8) {
	if i == 0 { // 000
		red.Low()
		green.Low()
		blue.Low()
	} else if i == 1 { // 001
		red.Low()
		green.Low()
		blue.High()
	} else if i == 2 { // 010
		red.Low()
		green.High()
		blue.Low()
	} else if i == 3 { // 011
		red.Low()
		green.High()
		blue.High()
	} else if i == 4 { // 100
		red.High()
		green.Low()
		blue.Low()
	} else if i == 5 { // 101
		red.High()
		green.Low()
		blue.High()
	} else if i == 6 { // 110
		red.High()
		green.High()
		blue.Low()
	} else if i == 7 { // 111
		red.High()
		green.High()
		blue.High()
	}
}

func main() {
	green.Configure(machine.PinConfig{Mode: machine.PinOutput})
	red.Configure(machine.PinConfig{Mode: machine.PinOutput})
	button.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	touch.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	blue.Configure(machine.PinConfig{Mode: machine.PinOutput})

	var isFirstButtonPressed bool
	var isSecondButtonPressed bool
	var isThirdButtonPressed bool
	var isFourthButtonPressed bool

	var input1 uint8
	var input2 uint8

	for {
		if !isFirstButtonPressed {
			// initialize all three leds
			red.High()
			green.High()
			blue.High()
		}

		if !isFirstButtonPressed && button.Get() {
			isFirstButtonPressed = true
			red.Low()
			green.Low()
			blue.Low()
			time.Sleep(time.Millisecond * 500)
		}

		if isFirstButtonPressed && !isSecondButtonPressed {
			if touch.Get() {
				if input1 < 3 {
					input1++
				}
				time.Sleep(time.Millisecond * 500)
			}

			if button.Get() {
				isSecondButtonPressed = true
				numberToLED(input1)
				time.Sleep(time.Millisecond * 500)
			}
		}

		if isSecondButtonPressed && !isThirdButtonPressed {
			if touch.Get() {
				if input2 < 3 {
					input2++
				}
				time.Sleep(time.Millisecond * 500)
			}

			if button.Get() {
				isThirdButtonPressed = true
				numberToLED(input2)
				time.Sleep(time.Millisecond * 500)
			}
		}

		if isThirdButtonPressed && !isFourthButtonPressed {
			if button.Get() {
				isFourthButtonPressed = true
				numberToLED(input1 + input2)
				time.Sleep(time.Millisecond * 500)
			}
		}

		if isFourthButtonPressed {
			if button.Get() {
				isFirstButtonPressed = false
				isSecondButtonPressed = false
				isThirdButtonPressed = false
				isFourthButtonPressed = false
				input1 = 0
				input2 = 0
				time.Sleep(time.Millisecond * 500)
			}
		}
	}
}
