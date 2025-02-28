// message type
// message identificator (id)

package main

import (
	"fmt"
	"log"
	"net"
	okaiparsetools "okai/common/okai-parse-tools"
	okaiparser "okai/common/okai-parser"
)

const (
	BUFF_SIZE = 5000
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
		msg := string(buff)
		pck := okaiparsetools.CutPacket(msg, "$")
		params := okaiparsetools.SplitParams(pck, ",")

		// fmt.Println("-------------------------------")
		// fmt.Println("Cutted pck:", pck)
		// fmt.Println("Splitted params:", params)

		parsed, _ := okaiparser.ParseParams(params)
		fmt.Println(parsed)
	}
}

func abortTCP(conn *Connection) {
	// do delete logic here
	// if connections[conn.Device.IMEI] != nil {
	// 	delete(connections, conn.Device.IMEI)
	// }
}

func oneStepParse(pck string) {
	params := okaiparsetools.SplitParams(pck, ",")
	data, _ := okaiparser.ParseParams(params)
	fmt.Println(data)
}

func main() {
	connections = make(map[string]*Connection)
	addr := fmt.Sprintf(":%d", TCP_PORT)

	serve, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalln("Startup serve error:", err.Error())
	}

	log.Println("Server started:", serve.Addr().Network())

	// oneStepParse("+RESP:GTFRI,OK043A,868070043228349,zk200,,,,,0,0000000000000000000,,,,,,,,,0250,0099,04E9,08C41A65,26&99,2,41,0,52322,4022,87,0,,0,,0.0&0.00&42.50&263.13&1&1&0&0000000D0000000D011A0000&000000D10500000000060000&00000000&0&02641C1B1A1AFFFFFF7D&1&1&00000000000000,85,20250228065940,00A6$0250228065921,009D$")

	for {
		conn, err := serve.Accept()

		if err != nil {
			log.Fatalln("accept connection error:", err.Error())
		}

		fmt.Printf("Received new connection:\n%v\n%v\n\r", conn.RemoteAddr().String(), conn.LocalAddr().String())

		go handleServe(conn)
	}
}
