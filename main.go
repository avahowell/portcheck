package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/johnathanhowell/reqlimit"
	"github.com/julienschmidt/httprouter"
)

type PortStatus struct {
	Open  bool
	Error string
}

// CheckPortHandler checks if the supplied address is connectable.
func CheckPortHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	st := PortStatus{}

	conn, err := net.DialTimeout("tcp", ps.ByName("addr"), time.Second*10)
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

	// limit requests to prevent abuse
	limitedRouter := reqlimit.New(router, 10, time.Minute)

	listenAddr := flag.String("bind", ":8080", "address to listen on")
	flag.Parse()

	log.Println("portcheck is listening on", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, limitedRouter))
}
