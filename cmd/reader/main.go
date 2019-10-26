package main

import (
	"encoding/binary"
	"fmt"
	"github.com/JHeimbach/nfc-cash-system/pkg/nfcreader"
	"github.com/fuzxxl/nfc/dev/nfc"
	"log"
	"os"
)

func main() {
	fmt.Println(nfc.ListDevices())

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	//var reciever = flag.String("receiver", "0.0.0.0:4433", "send cardId to this address")
	if isArgsLongEnough(1) {
		return fmt.Errorf("no command arg found")
	}
	switch os.Args[1] {
	case "poll":
		if isArgsLongEnough(2) {
			return fmt.Errorf("command poll: need devicename")
		}
		deviceName := os.Args[2]
		dev, err := nfcreader.OpenDevice(deviceName)
		if err != nil {
			return err
		}
		defer dev.Close()

		listenChan := dev.Listen()
		for {
			uidBytes := <-listenChan
			uidStr := binary.BigEndian.Uint32(uidBytes)
			fmt.Println(fmt.Sprint(uidStr))
		}
	}
	return nil
}
func isArgsLongEnough(minLength int) bool {
	return len(os.Args) < (minLength + 1) // + 1 for program name
}
