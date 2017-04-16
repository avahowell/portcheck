package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type PortStatus struct {
	Open  bool
	Error string
}

var timeout = time.Second * 5

func CheckPortHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	st := PortStatus{}

	conn, err := net.DialTimeout("tcp", ps.ByName("addr"), timeout)
	if err != nil {
		st.Open = false
		st.Error = err.Error()
	} else {
		conn.Close()
		st.Open = true
	}

	err = json.NewEncoder(w).Encode(st)
	if err != nil {
		log.Println("error encoding status: ", err)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/:addr", CheckPortHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
