package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// SrvConfig defines the far-end server, and its command and payload ports
type SrvConfig struct {
	Host    string
	RPCPort string
	Count   uint64
	Repeat  bool
}

// Command controls the type of function that TCPClient should perform
type Command struct {
	Name string
	Cfg  SrvConfig
}

type BitRate uint64

// Returns bitrate in Mega-bits per second
func (b BitRate) Mbps() float32 {
	return float32(b) / float32(1000000)
}

// Returns bitrate in Mega-bytes per second
func (b BitRate) MBps() float32 {
	return float32(b) / float32(8*1000000)
}

// Returns bitrate in Kilo-bits per second
func (b BitRate) Kbps() float32 {
	return float32(b) / float32(1000)
}

// Returns bitrate in Kilo-bytes per second
func (b BitRate) KBps() float32 {
	return float32(b) / float32(8*1000)
}

func (b BitRate) String() string {
	return fmt.Sprintf("%d", b)
}

// Stats is type of measurement that TCPClient reports on its stats channel.
type Stats struct {
	Stat string
	Type string
	Rate BitRate
}

type JSONStats struct {
	Stat string
	Type string
	Rate float32
}

// CCmdHandler is the receiver type for handling TCPClient control request
type CCmdHandler struct {
	CmdCh chan Command
}

// CStatHandler is the reciever type for handling TCPClient stats requests
type CStatHandler struct {
	StatCh chan Stats
}

// This handler parses the form from the user and initiates a TCPClient
// measurement.
func (c *CCmdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		raddr  string
		rport  int
		pktt   string
		tstt   string
		txsize int
		txmult string
		txcont string
	)
	params := map[string]interface{}{
		"raddr":  &raddr,
		"rport":  &rport,
		"pktt":   &pktt,
		"tstt":   &tstt,
		"txsize": &txsize,
		"txmult": &txmult,
		"txcont": &txcont,
	}
	Mult := map[string]uint64{
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
	}

	getformparams(r, params)
	trace.Printf("|CMD|%s|%s|\n", tstt, raddr)
	log.Println("CMD: ", raddr, rport, pktt, tstt, txsize, txmult, txcont != "")

	cmd := Command{
		Name: tstt,
		Cfg: SrvConfig{
			Host:    raddr,
			RPCPort: fmt.Sprint(rport),
			Count:   uint64(txsize) * Mult[txmult],
			Repeat:  txcont != "",
		},
	}
	c.CmdCh <- cmd
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(""))
}

// parse and store form parameters in the map that's passed in
func getformparams(r *http.Request, params map[string]interface{}) {
	for i, x := range params {
		if v := r.FormValue(i); v != "" {
			fmt.Sscan(v, x)
		}
	}
}

// This handler deals with GET requests for TCPClient measurement results.
// It returns the measurements since last snapshot in json format.
func (s *CStatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var jst JSONStats
	w.Header().Set("Content-Type", "text/plain")
	st, ok := <-s.StatCh
	if !ok {
		jst = JSONStats{Stat: "Error"}
	} else {
		jst = JSONStats{st.Stat, st.Type, st.Rate.Mbps()}
	}
	je := json.NewEncoder(w)
	je.Encode(jst)
}

// WebUI is an http server that provides an html UI to the user, annoucing itself at address
// that is passed in. It handles requests for starting and stopping of the load testing
// client and reporting of data.
func WebUI(addr string, cch chan Command, sch chan Stats) {
	cl := &CCmdHandler{cch}
	st := &CStatHandler{sch}
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.Handle("/cmd", cl)
	http.Handle("/stats", st)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: " + err.Error())
	}
}
