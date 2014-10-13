package fgms

import(
	"fmt"
	"log"
	"net"
	"time"
)

type crossfeed struct {
	Chan chan []byte
	Hosts map[string]*UDP_Conn
	Failed int
	Sent int
	MT_Failed int
	MT_Sent int
}

var CrossFeed *crossfeed

// auto initialize
func init() {

	CrossFeed = new(crossfeed)
	CrossFeed.Chan = make(chan []byte)
	CrossFeed.Hosts = make(map[string]*UDP_Conn)

	go CrossFeed.StartCheckTimer()
	go CrossFeed.Listen()
}

// add a host
func (me *crossfeed) Add(addr string, port int){

	conn := NewUDPConn(addr, port)
	me.Hosts[conn.Url] = conn

	go me.InitializeConn(conn)
}

// Attempt to connect and setup a Connection
// Should dns fail, address not exist or not able to connect
// the connection witll be marked as conn.Active = false with conn.LastError
func (me *crossfeed) InitializeConn( conn *UDP_Conn){

	if conn.Active {
		return // we need to define when its inactive, eg connection dropped
	}

	// get ip address
	udp_addr, err := net.ResolveUDPAddr("udp4", conn.Url)
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
	log.Println("Crossfeed: Connected to ", conn.Url )
}

// Starts a timer to check servers that are down  every 60 secs
// this is started in Loop() as a goroutine go me.StartCrossfeedCheckTimer()
func (me *crossfeed) StartCheckTimer(){

	ticker := time.NewTicker(time.Second * 60) // TODO roll back
	go func() {
		for _ = range ticker.C {
			for _, conn := range me.Hosts {
				if conn.Active == false {
					log.Println("> Attempt Reconnect Crossfeed = ", conn.Url )
					go me.InitializeConn(conn)
				}
			}
		}
	}()
}

// Listen for xdr packets from channel, and send to xrossfeeds
func (me *crossfeed) Listen(){
	fmt.Println("Crossfeed: Listening")
	for {
		select {
		case xdr_bytes := <- me.Chan:
			// Got data from channel
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
}

