package nfcreader

import (
	"fmt"
	"github.com/fuzxxl/nfc/dev/nfc"
	"io"
	"time"
)

type passiveTargetLister interface {
	io.Closer
	InitiatorListPassiveTargets(m nfc.Modulation) ([]nfc.Target, error)
}

type Device struct {
	device               passiveTargetLister
	PollingTimeOut       time.Duration
	AllowMultipleTargets bool
	Modulation           nfc.Modulation
}

func (d *Device) Close() error {
	return d.device.Close()
}

func (d *Device) Listen() <-chan []byte {
	send := make(chan []byte, 0)
	go func(send chan []byte) {
		if err := d.listenForTargets(send); err != nil {
			close(send)
		}
	}(send)
	return send
}

func (d Device) listenForTargets(send chan []byte) error {
	for {
		targets, err := d.device.InitiatorListPassiveTargets(d.Modulation)
		if err != nil {
			return fmt.Errorf("could not listen for targets: %v", err)
		}

		if len(targets) == 0 {
			time.Sleep(d.PollingTimeOut)
			continue
		}

		if d.AllowMultipleTargets {
			for _, t := range targets {
				sendTargetUid(t, send)
			}
			continue
		}

		sendTargetUid(targets[0], send)
	}
}

func sendTargetUid(target nfc.Target, send chan []byte) {
	card, ok := target.(*nfc.ISO14443aTarget)
	if ok {
		send <- card.UID[:card.UIDLen-1]
	}
}

func OpenDevice(name string) (*Device, error) {
	pnd, err := nfc.Open(name)

	if err != nil {
		return nil, fmt.Errorf("could not open device %q: %v", name, err)
	}

	return &Device{
		device:         pnd,
		PollingTimeOut: 100 * time.Millisecond,
		Modulation:     nfc.Modulation{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
	}, nil
}
