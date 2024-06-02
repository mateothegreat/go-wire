package receiver

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/mateothegreat/go-wire/test"
)

func TestReceiver(t *testing.T) {
	ch := make(chan *test.Image)
	go NewTCPReceiver("127.0.0.1", 15000, ch)

	go func() {
		for {
			select {
			case image := <-ch:
				log.Printf("Received image from %s, len=%d", image.Camera, len(image.Data))
				ioutil.WriteFile("output.jpg", image.Data, 0644)

			}
		}
	}()

	select {}
}
