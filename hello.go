package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var conns []*websocket.Conn

type User struct {
	Id   int
	Name string
	Age  int
}

type postBodySwitchMode struct {
	Address string `json: address`
	Mode    string `json: mode`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Format command for settings a input to multiple outputs
func inputToOutputs(input int, outputs ...int) string {
	command := fmt.Sprintf("x%dAV", input)
	outputStrs := make([]string, len(outputs))
	for i, output := range outputs {
		outputStrs[i] = fmt.Sprintf("x%d", output)
	}
	return command + strings.Join(outputStrs, ",")
}

func sendSignal(address string, commands []string) {
	logNPush(fmt.Sprintf("Sending signals... -> %s\n", address))

	logNPush(fmt.Sprintf("Dialing %s...", address))

	conn, err := net.Dial("tcp", address)
	if nil != err {
		logNPush(err.Error())
		return
	}

	defer conn.Close()

	readBuffer := bufio.NewReader(conn)

	for _, command := range commands {
		fmt.Fprintln(conn, command)
		response, err := readBuffer.ReadString('\n')
		if err != nil {
			logNPush(err.Error())
		} else {
			logNPush(response)
		}
	}
}

func logNPush(line string) {
	log.Println(line)
	pushLine(line)
}

func pushLine(line string) {
	for _, conn := range conns {
		err := conn.WriteMessage(websocket.TextMessage, []byte(line))
		if nil != err {
			log.Println(err)
		}
	}
}

func switchMode(mode, address string) {
	log.Println("Switching mode...", mode, address)
	if mode == "apple-tv" {
		sendSignal(address, []string{
			inputToOutputs(2, 1, 2),
		})
	} else if mode == "imac" {
		sendSignal(address, []string{
			inputToOutputs(1, 1, 2),
		})
	} else if mode == "chromecast" {
		sendSignal(address, []string{
			inputToOutputs(3, 1, 2),
		})
	}
}

func toJson(data interface{}) string {
	json, _ := json.MarshalIndent(data, "", "  ")
	return string(json)
}

func main() {
	log.Println("Starting...")
	runtime.GOMAXPROCS(runtime.NumCPU())

	r := mux.NewRouter()

	r.HandleFunc("/socket", func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println(err)
			return
		}

		conns = append(conns, conn)
	})

	r.Methods("POST").Path("/switch-mode").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var post postBodySwitchMode
		dec := json.NewDecoder(req.Body)
		err := dec.Decode(&post)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switchMode(post.Mode, post.Address)
	})

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(rice.MustFindBox("public").HTTPBox()),
	)
	n.UseHandler(r)
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = ":3000"
	}
	n.Run(listen)
}
