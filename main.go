package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
)

var SERVER_ADDRESS = ":5000"
var SIZE_OF_FILE = 8192

type FileServer struct {
}

func (fs *FileServer) start() {
	listener, err := net.Listen("tcp", SERVER_ADDRESS)
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal("Failed to Close server connection!\n err: ", err)
		}
	}(listener)
	if err != nil {
		log.Fatal("Failed to Start tcp server!\n err: ", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Failed to Accept server connection!\n err: ", err)
		}
		go fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	buffer := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		numBytes, err := io.CopyN(buffer, conn, size)
		if err != nil {
			log.Fatal("Failed to Read Buffer from connection!\n err: ", err)
		}

		log.Println(buffer.Bytes())
		log.Println("Received ", numBytes, " bytes over the network")
	}
}

func sendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		log.Fatal("Failed to send File\n err: ", err)
	}

	conn, err := net.Dial("tcp", SERVER_ADDRESS)
	if err != nil {
		log.Fatal("Failed to Dial connection\n err: ", err)
	}

	binary.Write(conn, binary.LittleEndian, int64(size))
	numBytes, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		log.Fatal("Failed to Write File\n err: ", err)
	}
	log.Printf("Written %d bytes over the network!", numBytes)

	return nil
}

func main() {
	go func() {
		time.Sleep(1 * time.Second)
		sendFile(SIZE_OF_FILE)
	}()
	server := new(FileServer)
	server.start()
}
