package sender

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"net"
	"sync"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

var (
	conn     net.Conn
	connLock sync.Mutex
)

func NewTCPSender(addr string) (chan bool, chan struct{}) {
	connected := make(chan bool)
	closeChan := make(chan struct{})
	reconnect := true

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			dialer := net.Dialer{}
			c, err := dialer.DialContext(ctx, "tcp", addr)
			if err != nil {
				log.Printf("Failed to connect to server: %v", err)
				time.Sleep(500 * time.Millisecond)
				continue
			}

			log.Printf("Connected to server")
			connected <- true
			connLock.Lock()
			conn = c
			connLock.Unlock()

			go func() {
				select {
				case <-closeChan:
					reconnect = false
					connLock.Lock()
					conn.Close()
					connLock.Unlock()
					connected <- false
					return
				}
			}()

			err = monitor(conn)
			if err != nil && reconnect {
				log.Printf("connection lost: %v", err)
			}

			conn.Close()
			connected <- false

			if !reconnect {
				break
			}

			log.Println("NewTCPSender connection closed... reconnecting")

			time.Sleep(500 * time.Millisecond)
		}
	}()

	return connected, closeChan
}

func Send(i any) error {
	var buf bytes.Buffer
	encoder := msgpack.NewEncoder(&buf)
	err := encoder.Encode(i)
	if err != nil {
		return err
	}

	data := buf.Bytes()
	length := int32(len(data))
	log.Printf("Sending data length: %d", length)
	err = binary.Write(conn, binary.LittleEndian, length)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func monitor(conn net.Conn) error {
	for {
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

		buffer := make([]byte, 1)
		_, err := conn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			return err
		}

		conn.SetReadDeadline(time.Time{})
	}
}
