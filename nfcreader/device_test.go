package nfcreader

import (
	"errors"
	"github.com/fuzxxl/nfc/dev/nfc"
	"github.com/google/go-cmp/cmp"
	isPkg "github.com/matryer/is"
	"testing"
	"time"
)

var targetId = [10]byte{0, 1, 2, 3, 4, 5, 6}
var targetIdSlice = targetId[:]

type mockReader struct {
	isClosed          bool
	concurrentTargets int
	createTargets     func(nfc.Modulation, int) ([]nfc.Target, error)
}

func (m *mockReader) Close() error {
	m.isClosed = true
	return nil
}

func (m *mockReader) InitiatorListPassiveTargets(mod nfc.Modulation) ([]nfc.Target, error) {
	return m.createTargets(mod, m.concurrentTargets)
}

func TestDevice_Listen(t *testing.T) {
	reader := &mockReader{
		createTargets: func(modulation nfc.Modulation, targetNum int) ([]nfc.Target, error) {
			targets := make([]nfc.Target, 0, targetNum)
			id := targetId
			for i := 0; i < targetNum; i++ {
				id[0] = byte(i)
				targets = append(targets, &nfc.ISO14443aTarget{
					UID:    id,
					UIDLen: len(targetId),
				})
			}

			return targets, nil
		},
	}
	t.Run("empty targets", func(t *testing.T) {
		is := isPkg.New(t)
		dev := &Device{device: reader}

		recv := makeChannelAndListen(t, dev)
		is.NoErr(dev.LastErr)

		got, ok := targetFromChannel(t, recv)
		if ok {
			t.Errorf("got target, did not expect target %v", got)
		}
	})
	t.Run("target recieved", func(t *testing.T) {
		is := isPkg.New(t)
		reader.concurrentTargets = 1
		defer func() {
			reader.concurrentTargets = 0
		}()
		dev := &Device{device: reader}

		recv := makeChannelAndListen(t, dev)
		is.NoErr(dev.LastErr)

		got, ok := targetFromChannel(t, recv)
		if !ok {
			t.Errorf("reciever timed out")
		}
		if cmp.Equal(got, targetIdSlice) {
			t.Errorf("got %v, expected %v", got, targetIdSlice)
		}
	})
	t.Run("multiple targets concurrently without AllowMultipleTargets", func(t *testing.T) {
		is := isPkg.New(t)
		reader.concurrentTargets = 2
		defer func() {
			reader.concurrentTargets = 0
		}()
		dev := &Device{device: reader}

		recv := makeChannelAndListen(t, dev)
		is.NoErr(dev.LastErr)

		got, ok := targetFromChannel(t, recv)
		if !ok {
			t.Errorf("reciever timed out")
		}

		if cmp.Equal(got, targetIdSlice) {
			t.Errorf("got %v, expected %v", got, targetIdSlice)
		}
	})
	t.Run("multiple targets concurrently with AllowMultipleTargets", func(t *testing.T) {
		is := isPkg.New(t)
		reader.concurrentTargets = 2
		defer func() {
			reader.concurrentTargets = 0
		}()
		dev := &Device{
			device:               reader,
			AllowMultipleTargets: true,
		}

		recv := makeChannelAndListen(t, dev)
		is.NoErr(dev.LastErr)

		for i := 0; i < 2; i++ {
			got, _ := targetFromChannel(t, recv)
			if got[0] != byte(i) {
				t.Errorf("got %v, expected %v", got[0], byte(i))
			}
		}
	})
	t.Run("multiple targets consecutively", func(t *testing.T) {
		is := isPkg.New(t)
		reader.concurrentTargets = 1
		defer func() {
			reader.concurrentTargets = 0
		}()
		dev := &Device{device: reader}
		recv := makeChannelAndListen(t, dev)
		is.NoErr(dev.LastErr)

		targetsRecv := 0
		for i := 0; i < 2; i++ {
			_, ok := targetFromChannel(t, recv)
			if ok {
				targetsRecv++
			}
		}
		if targetsRecv != 2 {
			t.Errorf("got %d targets, wanted 2", targetsRecv)
		}
	})
	t.Run("targets lister returns error", func(t *testing.T) {
		reader := &mockReader{
			concurrentTargets: 1,
			createTargets: func(modulation nfc.Modulation, targetNum int) ([]nfc.Target, error) {
				return nil, errors.New("test error")
			},
		}
		dev := &Device{
			device: reader,
		}
		recv := makeChannelAndListen(t, dev)

		select {
		case _, ok := <-recv:
			if !dev.HasError() {
				t.Errorf("device should have error")
			}
			if ok {
				t.Errorf("channel got not closed")
			}
		case <-time.After(10 * time.Minute):
			t.Errorf("timeout reciever")
		}
	})
}

func TestDevice_Close(t *testing.T) {
	reader := &mockReader{}
	dev := Device{
		device: reader,
	}
	_ = dev.Close()
	if !reader.isClosed {
		t.Errorf("passiveTargetLister close() was not called")
	}
}

func targetFromChannel(t *testing.T, recv <-chan []byte) ([]byte, bool) {
	t.Helper()
	select {
	case target := <-recv:
		return target, true
	case <-time.After(10 * time.Millisecond):
		return nil, false
	}
}

func makeChannelAndListen(t *testing.T, dev *Device) chan []byte {
	t.Helper()
	recv := make(chan []byte)
	go func(chan []byte) {
		dev.ListenForCardUids(recv)
	}(recv)

	return recv
}
