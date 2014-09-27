package main

import (
	"encoding/json"
	"github.com/GeertJohan/go.rice"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"runtime"
)

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

func switchMode(mode, address string) {
	log.Println("Switching mode...", mode, address)
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
	})

	r.Methods("GET", "HEAD").Path("/data").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	})

	r.Methods("GET", "HEAD").Path("/data/{id}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	})

	r.Methods("DELETE").Path("/data").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
