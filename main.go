package main

import (
	"fmt"
	"log"
	"net"
)

const (
	BUFF_SIZE = 4096
	TCP_PORT  = 30331
)

type Connection struct {
	Device any
	Conn   *net.Conn
}

var connections map[string]*Connection

func handleServe(conn net.Conn) {
	connection := &Connection{
		Conn:   &conn,
		Device: 1,
	}

	buff := make([]byte, BUFF_SIZE)

	defer func() {
		abortTCP(connection)
		conn.Close()
	}()

	for {
		_, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Received data err:", err.Error())
			break
		}
		fmt.Println("Received message from OKAI:", string(buff))
	}
}

func abortTCP(conn *Connection) {
	// do delete logic here
	// if connections[conn.Device.IMEI] != nil {
	// 	delete(connections, conn.Device.IMEI)
	// }
}

func main() {
	connections = make(map[string]*Connection)
	addr := fmt.Sprintf(":%d", TCP_PORT)

	serve, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalln("Startup serve error:", err.Error())
	}

	log.Println("Server started:", serve.Addr().Network())

	for {
		conn, err := serve.Accept()

		if err != nil {
			log.Fatalln("accept connection error:", err.Error())
		}

		fmt.Printf("Received new connection:\n%v\n%v\n\r", conn.RemoteAddr().String(), conn.LocalAddr().String())

		go handleServe(conn)
	}
}
