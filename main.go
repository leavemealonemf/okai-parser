package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	okaiparsetools "okai/common/okai-parse-tools"
	okaiparser "okai/common/okai-parser"
	"okai/common/utils"
	mg "okai/db/mg"
	"okai/rabbit"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	BUFF_SIZE = 5000
	TCP_PORT  = 30331
)

type Connection struct {
	IMEI       string
	Conn       net.Conn
	TotalCount string
}

type QueueCmd struct {
	CMD            string `json:"Cmd"`
	ImeiWithPrefix string `json:"IMEI"`
}

type ReceivedCommand struct {
	ServerTime    int64         `json:"_ts" bson:"_ts"`
	CompletedTime int64         `json:"_ct" bson:"_ct"`
	CMD           string        `json:"hex_origin" bson:"hex_origin,omitempty"`
	Token         string        `json:"token" bson:"token,omitempty"`
	Status        string        `json:"status" bson:"status,omitempty"`
	IMEI          string        `json:"dvce_imei" bson:"dvce_imei,omitempty"`
	CMDInfo       interface{}   `json:"cmd_info" bson:"cmd_info"`
	QueueD        amqp.Delivery `json:"-" bson:"-"`
	ExecChannel   chan bool     `json:"-" bson:"-"`
}

var rbtConn *amqp.Connection
var rbtCh *amqp.Channel

var connections map[string]*Connection
var receivedCommands map[string]*ReceivedCommand

func handleServe(conn net.Conn) {
	connection := &Connection{
		Conn: conn,
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

		// fmt.Printf("pId: %v\nPacket: %v\n", pId, parsed)

		if !authorized && pId == "GTNCN" {
			imei := parsed["imei"].(string)
			connection.IMEI = imei
			connections[imei] = connection
			tc := parsed["totalCount"].(string)
			connection.TotalCount = tc[0:4]
			authorized = true
			fmt.Println("succesfully authorized")
			continue
		}

		if !authorized {
			continue
		}

		if pType == "+ACK" {
			if pId == "GTECC" || pId == "GTRTO" {
				cmdID := parsed["cmdID"].(string)
				fmt.Println(cmdID)
				receivedCommand := receivedCommands[cmdID]
				if receivedCommand != nil {
					receivedCommand.ExecChannel <- true
				}
			}
			continue
		}

		// test case gnss info
		if pType == "+RESP" && pId == "GTFRI" {
			fmt.Println("RAW GTFRI PACKET:", string(buff))
		}

		// heartbeat handshake
		if pType == "+ACK" && pId == "GTHBD" {
			protoVer := params[1]
			totalCount := params[6]
			cmd := fmt.Sprintf("+SACK:GTHBD,%s,%s", protoVer, totalCount)
			conn.Write([]byte(cmd))
			fmt.Println("send heartbeat ack:", cmd)
		}

		if len(parsed) > 0 {
			tc := parsed["totalCount"].(string)
			connection.TotalCount = tc[0:4]
			err = insertOneScooter(ctx, parsed)
			if err == nil {
				fmt.Println("insert successfully")
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}

func insertOneScooter(ctx context.Context, dvce map[string]interface{}) error {
	lat := dvce["lat"].(string)
	lon := dvce["lon"].(string)
	if len(lat) == 0 && len(lon) == 0 {
		var latest bson.M
		mg.FindOneWithOpts(
			ctx, scooterColl,
			bson.D{{Key: "imei", Value: dvce["imei"]}},
			options.FindOne().SetSort(bson.D{{Key: "_ts", Value: -1}}),
		).Decode(&latest)
		if latest != nil {
			dvce["lat"] = latest["lat"]
			dvce["lon"] = latest["lon"]
		}
	}

	_, e := scooterColl.InsertOne(ctx, dvce)

	marshal, _ := json.Marshal(dvce)
	publishPacket(marshal)

	if e != nil {
		return e
	}

	return nil
}

func abortTCP(conn *Connection) {
	if connections[conn.IMEI] != nil {
		delete(connections, conn.IMEI)
	}

	filter := bson.D{
		{Key: "imei", Value: conn.IMEI},
	}

	update := bson.D{
		{"$set", bson.D{
			{Key: "online", Value: false},
		}},
	}

	opts := options.FindOneAndUpdate().SetSort(bson.D{{Key: "_ts", Value: -1}})
	res := mg.UpdOneScooter(ctx, scooterColl, filter, update, opts)

	var dvce bson.M
	err := res.Decode(&dvce)
	if err != nil {
		fmt.Println("[abort tcp conn] failed to decode result mongo")
		return
	}
	dvce["online"] = false
	fmt.Println("online info after disconnect", dvce["online"])
	j, err := json.Marshal(&dvce)
	if err != nil {
		fmt.Println("[abort tcp conn] failed to marshal result")
		return
	}

	publishPacket(j)
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

var commands map[string]map[string]string

func initCommands() {
	commands, _ = utils.LoadJSON[map[string]map[string]string]("commands.json")
	fmt.Println(commands)
}

func HTTPCommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "failed to read bytes", 400)
		}

		var cmdS QueueCmd
		err = json.Unmarshal(body, &cmdS)
		if err != nil {
			http.Error(w, "Failed to parse JSON response", http.StatusInternalServerError)
			return
		}

		resStr := strings.Split(cmdS.ImeiWithPrefix, ":")
		cmd := cmdS.CMD
		imei := resStr[1]

		c := connections[imei]

		if c != nil {
			cmdInfo := commands[cmd]

			if cmdInfo == nil {
				w.WriteHeader(404)
				w.Write([]byte("this command does not exist"))
				return
			}

			token := c.TotalCount

			bCommand := okaiparser.CommandBuilder(cmdInfo, token)
			cmdChan := make(chan bool)

			recievedCmd := &ReceivedCommand{
				ServerTime:  time.Now().UnixMicro(),
				CMD:         bCommand,
				Token:       token,
				Status:      "pending",
				IMEI:        imei,
				CMDInfo:     commands[cmd],
				ExecChannel: cmdChan,
			}

			receivedCommands[token] = recievedCmd
			// mg.Insert(ctx, cmdsColl, recievedCmd)
			c.Conn.Write([]byte(bCommand))

			select {
			case success := <-cmdChan:
				if success {
					w.Write([]byte(fmt.Sprintf("Command %s executed successfully", cmd)))
					delete(receivedCommands, token)
				} else {
					http.Error(w, "Command execution failed", http.StatusInternalServerError)
					delete(receivedCommands, token)
				}
			case <-time.After(60 * time.Second):
				http.Error(w, "Command execution timed out", http.StatusRequestTimeout)
				delete(receivedCommands, token)
			}

		} else {
			w.WriteHeader(404)
			w.Write([]byte("device not connected"))
			return
		}
	} else {
		http.Error(w, "method not implemented", 404)
		return
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	initCommands()
	fmt.Println(commands["turnOn"])
}

func publishPacket(pkt []byte) {
	err := rbtCh.Publish(
		"",
		"packets",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        pkt,
		},
	)
	if err != nil {
		log.Println("failed to publish message", err.Error())
	}
}

func initHttp() {
	r := mux.NewRouter()

	r.HandleFunc("/cmd", HTTPCommandHandler)

	addr := fmt.Sprintf("0.0.0.0:%d", 8912)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	srv := &http.Server{
		Handler:      handler,
		Addr:         addr,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	fmt.Printf("http server up on :%d\n", 8912)
	log.Fatal(srv.ListenAndServe())
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

	receivedCommands = map[string]*ReceivedCommand{}

	go showConnections()
	go initHttp()

	rbtConn = rabbit.Conn()
	ch, err := rbtConn.Channel()
	rbtCh = ch

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer ch.Close()

	_, err = ch.QueueDeclare(
		"packets",
		false,
		false,
		false,
		false,
		nil,
	)
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
