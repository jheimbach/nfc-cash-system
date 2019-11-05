package nfcreader

import (
	"fmt"
	"github.com/fuzxxl/nfc/2.0/nfc"
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
	LastErr              error
}

func (d *Device) Close() error {
	return d.device.Close()
}

func (d *Device) ListenForCardUids(send chan []byte) {
	for {
		targets, err := d.device.InitiatorListPassiveTargets(d.Modulation)
		if err != nil {
			d.LastErr = fmt.Errorf("could not listen for targets: %v", err)
			close(send)
			break
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

func (d *Device) HasError() bool {
	return d.LastErr != nil
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
