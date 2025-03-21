package main

import (
	"context"
	"fmt"
	"log"
	"net"
	okaiparsetools "okai/common/okai-parse-tools"
	okaiparser "okai/common/okai-parser"
	"okai/common/utils"
	mg "okai/db/mg"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
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

		pck := okaiparsetools.CutPacket(msg, "$")
		params := okaiparsetools.SplitParams(pck, ",")

		pType, pId, parsed, _ := okaiparser.ParseParams(params)

		if !authorized && pId == "GTNCN" {
			imei := parsed["imei"].(string)
			connection.IMEI = imei
			connections[imei] = connection
			authorized = true
			fmt.Println("succesfully authorized")
			continue
		}

		if !authorized {
			continue
		}

		// heartbeat handshake
		if pType == "+ACK" && pId == "GTHBD" {
			protoVer := params[1]
			totalCount := params[6]
			cmd := fmt.Sprintf("+SACK:GTHBD,%s,%s", protoVer, totalCount)
			conn.Write([]byte(cmd))
			fmt.Println("send heartbeat ack:", cmd)
		}

		_, err = utils.JsonStringify(parsed)

		if err == nil {
			_, err = scooterColl.InsertOne(ctx, parsed)
			if err == nil {
				fmt.Println("insert successfully")
			} else {
				fmt.Println(err.Error())
			}
		}
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
	for {
		fmt.Printf("========== CONNECTIONS ===========\n")
		cLen := len(connections)
		if cLen > 0 {
			for _, v := range connections {
				fmt.Println(v.IMEI)
			}
		}
		fmt.Printf("==================================\n")
		time.Sleep(time.Second * 60)
	}
}

var scooterColl *mongo.Collection
var ctx = context.TODO()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	mongUsr, _ := os.LookupEnv("MONGO_USR")
	mongPass, _ := os.LookupEnv("MONGO_PASSWORD")
	mongHost, _ := os.LookupEnv("MONGO_HOST")
	mongPort, _ := os.LookupEnv("MONGO_PORT")

	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%s", mongUsr, mongPass, mongHost, mongPort)
	mgClient, err := mg.Connect(ctx, connStr)
	mg.Seed(mgClient, ctx)
	scooterColl = mgClient.Database("iot").Collection("okai_scooters")

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
