package receiver

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"

	"github.com/vmihailenco/msgpack/v5"
)

func NewTCPReceiver[T any](addr string, port int, ch chan<- T) error {
	log.Printf("Listening on port %d", port)

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: port,
	})
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Error accepting connection: %s", err.Error())
			continue
		}

		log.Printf("Accepted connection from %s", conn.RemoteAddr().String())

		go handleConnection(conn, ch)
	}
}

func handleConnection[T any](conn *net.TCPConn, ch chan<- T) {
	defer conn.Close()

	for {
		var length int32
		err := binary.Read(conn, binary.LittleEndian, &length)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by client")
			} else {
				log.Printf("Error reading length from connection: %s", err.Error())
			}
			return
		}

		log.Printf("Reading data of length: %d", length)
		data := make([]byte, length)
		_, err = io.ReadFull(conn, data)
		if err != nil {
			log.Printf("Error reading data from connection: %s", err.Error())
			return
		}

		log.Printf("Received %d bytes", length)

		var t T
		reader := bytes.NewReader(data)
		decoder := msgpack.NewDecoder(reader)
		err = decoder.Decode(&t) // Pass a pointer to t
		if err != nil {
			log.Printf("Error unmarshalling data: %s", err.Error())
			continue
		}

		ch <- t
	}
}
