package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	minidrone "github.com/hybridgroup/tinygo-minidrone"
	term "github.com/nsf/termbox-go"
	"tinygo.org/x/bluetooth"
)

var deviceAddress = connectAddress()

var (
	adapter = bluetooth.DefaultAdapter
	device  bluetooth.Device
	ch      = make(chan bluetooth.ScanResult, 1)
	drone   *minidrone.Minidrone
)

func main() {
	defer cleanup()

	fmt.Println("Enabling Bluetooth...")
	must("enable BLE interface", adapter.Enable())

	fmt.Println("Starting scan...")
	must("start scan", adapter.Scan(scanHandler))

	select {
	case result := <-ch:
		var err error
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		must("connect to device", err)
		fmt.Printf("Connected to %s\n", result.Address.String())
	case <-time.After(30 * time.Second):
		fmt.Println("Scan timeout. No matching device found.")
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Exiting...")
		// Perform any necessary cleanup here
		os.Exit(0)
	}()

	drone = minidrone.NewMinidrone(&device)
	must("start drone", drone.Start())

	fmt.Println("Initializing terminal...")
	must("initialize terminal", term.Init())

	go handleSignals()

	fmt.Println("Taking off in 3 seconds...")
	time.Sleep(3 * time.Second)
	drone.TakeOff()

	controlDrone()
}

func controlDrone() {
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				fmt.Println("Exiting...")
				return
			default:
				handleDroneCommand(ev.Ch)
			}
		case term.EventError:
			fmt.Printf("Terminal error: %v\n", ev.Err)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func handleDroneCommand(ch rune) {
	commands := map[rune]func(){
		'[': func() { fmt.Println("Takeoff..."); drone.TakeOff() },
		']': func() { fmt.Println("Land..."); drone.Land() },
		'w': func() { fmt.Println("Forward..."); drone.Forward(1) },
		's': func() { fmt.Println("Backward..."); drone.Backward(1) },
		'a': func() { fmt.Println("Left..."); drone.Left(1) },
		'd': func() { fmt.Println("Right..."); drone.Right(1) },
		'k': func() { fmt.Println("Down..."); drone.Down(1) },
		'i': func() { fmt.Println("Up..."); drone.Up(1) },
		'j': func() {
			fmt.Println("Spin counter clockwise...")
			drone.CounterClockwise(minidrone.ValidatePitch(20, 10))
		},
		'l': func() { fmt.Println("Spin clockwise..."); drone.Clockwise(minidrone.ValidatePitch(20, 10)) },
		't': func() { fmt.Println("Front flip..."); drone.FrontFlip() },
		'g': func() { fmt.Println("Back flip..."); drone.BackFlip() },
		'f': func() { fmt.Println("Left flip..."); drone.LeftFlip() },
		'h': func() { fmt.Println("Right flip..."); drone.RightFlip() },
	}

	if cmd, ok := commands[ch]; ok {
		cmd()
	} else {
		drone.Halt()
	}
}

func cleanup() {
	fmt.Println("Cleaning up...")
	if drone != nil {
		drone.Land()
		time.Sleep(2 * time.Second)
		drone.Halt()
	}

	device.Disconnect()

	if adapter != nil {
		adapter.StopScan()
	}
	term.Close()
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\nReceived termination signal")
	cleanup()
	os.Exit(0)
}

func scanHandler(a *bluetooth.Adapter, d bluetooth.ScanResult) {
	fmt.Printf("Device: %s, RSSI: %d, Name: %s\n", d.Address.String(), d.RSSI, d.LocalName())
	if d.Address.String() == deviceAddress {
		a.StopScan()
		ch <- d
	}
}

func must(action string, err error) {
	if err != nil {
		fmt.Printf("Failed to %s: %v\n", action, err)
		os.Exit(1)
	}
}

func connectAddress() string {
	if len(os.Args) < 2 {
		println("you must pass the Bluetooth address of the minidrone y0u want to connect to as the first argument")
		os.Exit(1)
	}
