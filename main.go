package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var conns []*websocket.Conn

// PostBodySwitchMode is data returned from the client.
// Mode is what input should be switched to.
// Address is the HDMI switchers IP.
type PostBodySwitchMode struct {
	Address string `json:"address"`
	Mode    string `json:"mode"`
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

	conn, err := net.DialTimeout("tcp", address, 7*time.Second)
	if nil != err {
		logNPush(err.Error())
		return
	}

	defer conn.Close()

	readBuffer := bufio.NewReader(conn)

	logNPush("// Reading some lines")

	for i := 0; i < 2; i++ {
		response, err := readBuffer.ReadString('\n')
		if nil != err {
			logNPush(err.Error())
		}

		logNPush(response)
	}

	logNPush("// Read some lines...")

	for _, command := range commands {
		logNPush("-> " + command)
		fmt.Fprintln(conn, command+"\r")

		// Discard two lines
		_, err := readBuffer.ReadString('\n')
		_, err = readBuffer.ReadString('\n')
		response, err := readBuffer.ReadString('\n')

		if err != nil {
			logNPush("<- " + err.Error())
		} else {
			logNPush("<- " + response)
		}
	}
}

var lineCh chan string

func logNPush(line string) {
	log.Println(line)
	lineCh <- line
}

func switchMode(mode, address string) {
	log.Println("Switching mode...", mode, address)
	input := 0

	switch mode {
	case "mac-mini":
		input = 1
	case "apple-tv":
		input = 2
	case "x":
		input = 3
	case "chromecast":
		input = 4
	}

	sendSignal(address, []string{
		inputToOutputs(input, 1, 2, 3, 4),
	})
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
		var post PostBodySwitchMode
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

	lineCh = make(chan string, 32)

	go func() {
		for line := range lineCh {
			for _, conn := range conns {
				err := conn.WriteMessage(websocket.TextMessage, []byte(line))
				if nil != err {
					log.Println(err)
				}
			}
		}
	}()

	n.Run(listen)
}
