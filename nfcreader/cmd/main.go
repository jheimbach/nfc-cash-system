package main

import (
	"encoding/binary"
	"fmt"
	"github.com/JHeimbach/nfc-cash-system/nfcreader"
	"github.com/fuzxxl/nfc/2.0/nfc"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	//var reciever = flag.String("receiver", "0.0.0.0:4433", "send cardId to this address")
	if !isArgsLongEnough(1) {
		return fmt.Errorf("no command arg found")
	}
	switch os.Args[1] {
	case "poll":
		err := pollingDevice()
		if err != nil {
			return err
		}
	case "list":
		if err := listDevices(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("command %s not valid", os.Args[1])
	}
	return nil
}

func pollingDevice() error {
	deviceName, err := selectDevice()
	if err != nil {
		return err
	}
	dev, err := nfcreader.OpenDevice(deviceName)
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	defer dev.Close()

	// open channel to send uids to
	listenChan := make(chan []byte)

	// listen for targets in goroutine
	go func(chan []byte) {
		dev.ListenForCardUids(listenChan)
	}(listenChan)

	fmt.Printf("device %q ready, start polling...\n", deviceName)
	for {
		uidBytes, open := <-listenChan
		if !open {
			fmt.Printf("device %q: polling stopped \n", deviceName)
			break
		}
		uidStr := binary.BigEndian.Uint32(uidBytes)
		// todo do something with uid
		fmt.Println(fmt.Sprint(uidStr))
	}

	if dev.HasError() {
		return dev.LastErr
	}

	return nil
}

func isArgsLongEnough(minLength int) bool {
	return len(os.Args) >= (minLength + 1) // + 1 for program name
}

func listDevices() error {
	devices, err := nfc.ListDevices()
	if err != nil {
		return err
	}

	if len(devices) == 0 {
		fmt.Println("no devices found")
		return nil
	}

	for key, dev := range devices {
		fmt.Printf("[%d] %s \n", key, dev)
	}

	return nil
}

func selectDevice() (string, error) {
	if isArgsLongEnough(2) {
		return os.Args[2], nil
	}
	fmt.Printf("no device in arguments found, please select one from this list:\n")
	err := listDevices()
	if err != nil {
		return "", err
	}

	var selection = 0
	fmt.Printf("device num:")
	_, err = fmt.Scanf("%d", &selection)
	if err != nil {
		return "", err
	}

	devices, err := nfc.ListDevices()
	if err != nil {
		return "", err
	}

	return devices[selection], nil
}
