package fgms

import(
	"fmt"
	"log"
	"net"
	"time"
)

type crossfeed struct {
	Chan chan []byte
	Hosts map[string]*crossfeed_host
	Failed int
	Sent int
	//MT_Failed int
	//MT_Sent int
}

type crossfeed_host struct{
	Addr string
	Active bool
	LastError string
	Sock net.Conn
}


var Crossfeed *crossfeed



// Initialise and setup the `CrossFeed`
func SetupCrossfeed() {

	Crossfeed = new(crossfeed)
	Crossfeed.Chan = make(chan []byte)
	Crossfeed.Hosts = make(map[string]*crossfeed_host)

	go Crossfeed.Listen()
	go Crossfeed.StartReconnectTimer()

}

// add a host
func (me *crossfeed) Add(addr string, port int){

	host := new(crossfeed_host)
	host.Addr = fmt.Sprintf("%s:%d", addr, port)
	me.Hosts[host.Addr] = host

	go me.InitializeConn(host)
}

// Attempt to connect and setup a Connection
// Should dns fail, address not exist or not able to connect
// the connection will be marked as Active = false with LastError
func (me *crossfeed) InitializeConn( conn *crossfeed_host){

	if conn.Active {
		return // todo - need to define when its inactive, eg connection dropped
	}

	// dns lookup
	udp_addr, err := net.ResolveUDPAddr("udp4", conn.Addr)
	if err != nil {
		conn.Active = false
		conn.LastError = "Could not resolve IP address"
		return
	}

	//= open socket
	var err_listen error
	conn.Sock, err_listen = net.Dial("udp4", udp_addr.String())
	if err_listen != nil {
		conn.Active = false
		conn.LastError = "Could not open UDP port"
		return
	}

	// all good
	conn.Active = true
	conn.LastError = ""
	log.Println("Crossfeed: Connected to ", conn.Addr )
}

// Starts a timer to reconnect to hosts that are down  every 60 secs
// this is started as a goroutine
func (me *crossfeed) StartReconnectTimer(){

	ticker := time.NewTicker(time.Second * 60) // TODO roll back
	//go func() {
		for _ = range ticker.C {
			for _, host := range me.Hosts {
				if host.Active == false {
					log.Println("> Attempt Reconnect Crossfeed = ", host.Addr )
					go me.InitializeConn(host)
				}
			}
		}
	//}()
}

// Listen for xdr packets from channel, and send to xrossfeeds
func (me *crossfeed) Listen(){
	fmt.Println("Crossfeed: Listening")
	for {

		xdr_bytes := <- me.Chan

		for _, cf := range me.Hosts {
			if cf.Active {
				_, err := cf.Sock.Write(xdr_bytes)
				if err != nil {
					fmt.Println("Crossfeed error", err)
					me.Failed++
				}else {
					me.Sent++
				}
			}
		}

	}
}

