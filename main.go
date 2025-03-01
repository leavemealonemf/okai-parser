package main

import (
	"fmt"
	"log"
	"net"
	okaiparsetools "okai/common/okai-parse-tools"
	okaiparser "okai/common/okai-parser"
	"okai/common/utils"
)

const (
	BUFF_SIZE = 5000
	TCP_PORT  = 30331
)

type Connection struct {
	IMEI string
	Conn *net.Conn
}

var connections map[string]*Connection

func handleServe(conn net.Conn) {
	connection := &Connection{
		Conn: &conn,
		IMEI: "",
	}

	buff := make([]byte, BUFF_SIZE)

	authorized := false

	defer func() {
		abortTCP(connection)
		conn.Close()
	}()

	for {
		if authorized && connections[connection.IMEI] == nil {
			break
		}

		_, err := conn.Read(buff)

		if err != nil {
			fmt.Println("Received data err:", err.Error())
			break
		}
		msg := string(buff)

		if msg[0] != '+' {
			break
		}

		// fmt.Println("new message:", msg)
		// fmt.Println("--------------------------")

		pck := okaiparsetools.CutPacket(msg, "$")
		params := okaiparsetools.SplitParams(pck, ",")

		_, pId, parsed, _ := okaiparser.ParseParams(params)
		if parsed != nil {
			if !authorized && pId == "GTNCN" {
				imei := parsed["imei"].(string)
				connection.IMEI = imei
				connections[imei] = connection
				authorized = true
			} else {
				fmt.Println("break connection...")
				break
			}
		}
		jsn, _ := utils.JsonStringify(parsed)
		fmt.Println(jsn)
	}
}

func abortTCP(conn *Connection) {
	if connections[conn.IMEI] != nil {
		delete(connections, conn.IMEI)
	}
}

func oneStepParse(pck string) {
	params := okaiparsetools.SplitParams(pck, ",")
	_, _, parsed, _ := okaiparser.ParseParams(params)
	if parsed != nil {
		fmt.Println(parsed["imei"])
	}
	jsn, _ := utils.JsonStringify(parsed)
	fmt.Println(jsn)
}

func showConnections() {
	fmt.Printf("========== CONNECTIONS ===========")
	cLen := len(connections)
	if cLen > 0 {
		for _, v := range connections {
			fmt.Println(v.IMEI)
		}
	}
	fmt.Printf("==================================")
}

func main() {
	connections = make(map[string]*Connection)
	addr := fmt.Sprintf(":%d", TCP_PORT)

	serve, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalln("Startup serve error:", err.Error())
	}

	log.Println("Server started:", serve.Addr().Network())

	go showConnections()
	// oneStepParse("+RESP:GTFRI,OK043A,868070043228349,zk200,,,,,0,0000000000000000000,,,,,,,,,0250,0099,04E9,08C41A65,26&99,2,41,0,52322,4022,87,0,,0,,0.0&0.00&42.50&263.13&1&1&0&0000000D0000000D011A0000&000000D10500000000060000&00000000&0&02641C1B1A1AFFFFFF7D&1&1&00000000000000,85,20250228065940,00A6$0250228065921,009D$")
	// oneStepParse("+RESP:GTNCN,OK043A,868070043228349,zk200,,,,,1,0000000000000000000,4,1.0,,218.2,,,20250301154818,,0250,0099,04E9,08C41A65,29&99,2,41,0,51521,4034,87,0,0,1,,0.0&0.00&0.00&263.22&0&0&0&0000000D0000000D011A0000&000000D10500000000060000&00000000&0&00000000000000000000&1&1&00000000000000,78,20250301154821,000B$")
	for {
		conn, err := serve.Accept()

		if err != nil {
			log.Fatalln("accept connection error:", err.Error())
		}

		fmt.Printf("Received new connection:\n%v\n%v\n\r", conn.RemoteAddr().String(), conn.LocalAddr().String())

		go handleServe(conn)
	}
}
