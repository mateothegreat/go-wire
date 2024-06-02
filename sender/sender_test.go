package sender

import (
	"os"
	"testing"
	"time"

	"github.com/mateothegreat/go-wire/test"
)

func TestSend(t *testing.T) {
	wire := &WireTCPSender{}
	connected, closeChan := wire.Connect("127.0.0.1:15000")

	select {
	case <-connected:
		f, err := os.ReadFile("../test/test.jpg")
		if err != nil {
			t.Errorf("Error reading image: %v", err)
		}

		image := test.Image{
			Camera: "test",
			Data:   f,
		}

		err = wire.Send(image)
		if err != nil {
			t.Errorf("Error sending image: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for connection")
	}

	closeChan <- struct{}{}

	select {
	case <-connected:
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for disconnection")
	}
}
