
package fgms

import(
	"fmt"
	"log"
	"net"
	"time"
)

import(
	//"github.com/davecgh/go-xdr/xdr"
	//"github.com/FreeFlightSim/go-fgms/flightgear"
	//"github.com/FreeFlightSim/go-fgms/message"
)

type crossfeed struct {
	Chan chan []byte
	Hosts map[string]*UDP_Conn //*net.UDPConn
	Failed int
	Sent int
	MT_Failed int
	MT_Sent int
}

var CrossFeed *crossfeed

func init() {
	CrossFeed = new(crossfeed)
	CrossFeed.Chan = make(chan []byte)
	CrossFeed.Hosts = make(map[string]*UDP_Conn)
	go CrossFeed.Listen()
	fmt.Println("init ###################")
}

func (me *crossfeed) Add(addr string, port int){
	log.Println("> Add Crossfeed ================= ", addr, port)

	conn := NewUDPConn(addr, port)
	me.Hosts[conn.Url] = conn

	go me.InitConn(conn)
}

// Attempt to connect and setup a Crossfeed Connection
// Should dns fail, address not exist or not able to connect
// the connection witll be marked as conn.Active = false with conn.LastError
func (me *crossfeed) InitConn( conn *UDP_Conn){

	log.Println("> InitSetupCrossfeed = ", conn.Url )

	if conn.Active {
		return // we need to define when its inactive, eg connection dropped
	}

	//= resolve with UDP address
	udp_addr, err := net.ResolveUDPAddr("udp4", conn.Url)
	if err != nil {
		conn.Active = false
		conn.LastError = "Could nto resolve IP address"
		log.Println("\tFAIL: Crossfeed to resolve UDP address:  ", conn.Url, err)
		return
	}

	//= open socket and listen
	var err_listen error
	conn.Sock, err_listen = net.Dial("udp4", udp_addr.String())
	if err_listen != nil {
		conn.Active = false
		conn.LastError = "Couldnt open UDP port"
		log.Println("\tFAIL: Crossfeed FAIL to Open:  ", conn.Url, udp_addr, err_listen)
		return
	}
	conn.Active = true
	conn.LastError = ""
	log.Println("\tOK:   Crossfeed Added -  ", conn.Url, udp_addr, err_listen)

} // FgServer::AddCrossfeed()


func (me *crossfeed) Listen(){
	fmt.Println("Listen---------------")
	for {
		select {
		case xdr_bytes := <- me.Chan:
			for _, cf := range me.Hosts {
				if cf.Active {
					//if (CurrentCrossfeed->Address.getIP() != SenderAddress.getIP()) ?? But same address different port ??
					_, err := cf.Sock.Write(xdr_bytes)
					if err != nil {
						fmt.Println(err)
						me.Failed++
					}else {
						me.Sent++
					}
				}
				//CurrentCrossfeed++;
			}
		}
	}
	fmt.Println("DONEEE")
}







func (me *FgServer) AddCrossfeed( host_name string, port int){

	log.Println("> Add Crossfeed = ", host_name, port)

	// Create object
	conn := NewUDPConn(host_name, port)
	me.Crossfeeds[conn.Url] = conn

	// Go check and setup
	go me.InitSetupCrossfeed(conn)

}

// Attempt to connect and setup a Crossfeed Connection
// Should dns fail, address not exist or not able to connect
// the connection witll be marked as conn.Active = false with conn.LastError
func (me *FgServer) InitSetupCrossfeed( conn *UDP_Conn){

	log.Println("> InitSetupCrossfeed = ", conn.Url )

	if conn.Active {
		return // we need to define when its inactive, eg connection dropped
	}

	//= resolve with UDP address
	udp_addr, err := net.ResolveUDPAddr("udp4", conn.Url)
	if err != nil {
		conn.Active = false
		conn.LastError = "Could nto resolve IP address"
		log.Println("\tFAIL: Crossfeed to resolve UDP address:  ", conn.Url, err)
		return
	}

	//= open socket and listen
	var err_listen error
	conn.Sock, err_listen = net.Dial("udp4", udp_addr.String())
	if err_listen != nil {
		conn.Active = false
		conn.LastError = "Couldnt open UDP port"
		log.Println("\tFAIL: Crossfeed FAIL to Open:  ", conn.Url, udp_addr, err_listen)
		return
	}
	conn.Active = true
	conn.LastError = ""
	log.Println("\tOK:   Crossfeed Added -  ", conn.Url, udp_addr, err_listen)

} // FgServer::AddCrossfeed()


//= Starts a timer to check servers that are down ie Active = false (currently every 60 secs)
// this is started in Loop() as a goroutine go me.StartCrossfeedCheckTimer()
func (me *FgServer) StartCrossfeedCheckTimer(){
	
	ticker := time.NewTicker(time.Millisecond * 60000)
    go func() {
		for _ = range ticker.C {	
			for _, conn := range me.Crossfeeds {			
				if conn.Active == false {
					log.Println("> Attempt Reconnect Crossfeed = ", conn.Url )
					go me.InitSetupCrossfeed(conn)
				}
			}
		}
	}()
}	





/*   Send message to all crossfeed servers.
	Crossfeed servers receive all traffic without condition,
	mainly used for testing and debugging, and crossfeed.fgx.ch
	http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/FgServer.cxx#line1154
*/
func (me *FgServer) SendToCrossfeed(xdr_bytes []byte, sender_address *net.UDPAddr){
	//return
	//T_MsgHdr*       MsgHdr;
	//uint32_t        MsgMagic;
	//int             sent;
	
	//MsgHdr    = (T_MsgHdr *) Msg;
	//MsgMagic  = MsgHdr->Magic;
	//MsgHdr->Magic = XDR_encode<uint32_t> (RELAY_MAGIC);
	
	// Not sure what is happening, but we create another payload, by unmarshalling Msg again ?
	//var MsgHdr flightgear.T_MsgHdr
	//_, err := xdr.Unmarshal(Msg, &MsgHdr)
	//if err != nil {
	//	fmt.Println("XDR Decode Error in SendToCrossfeed - Should never happen?", err)
	//	return
	//}
	//MsgHdr.Magic = message.RELAY_MAGIC
	
	//encoded, err := xdr.Marshal(MsgHdr)
	//if err != nil {
	//	fmt.Println("XDR Encode Error in SendToCrossfeed - Should never happen?", err)
	//	return
	//}
	//mT_RelayListIt CurrentCrossfeed = m_CrossfeedList.begin();
	//while (CurrentCrossfeed != m_CrossfeedList.end())
	//{

	for _, cf := range me.Crossfeeds {
		if cf.Active {
			//if (CurrentCrossfeed->Address.getIP() != SenderAddress.getIP()) ?? But same address different port ??
			_, err := cf.Sock.Write(xdr_bytes)
			if err != nil {
				fmt.Println(err)
				me.CrossFeedFailed++
			}else {
				me.CrossFeedSent++
			}
		}
		//CurrentCrossfeed++;
	}
	//fmt.Println("crossfeeds", me.Crossfeeds, me.CrossFeedSent, me.CrossFeedFailed)
	//MsgHdr->Magic = MsgMagic;  // restore the magic value ? umm not used now ?
} // FgServer::SendToCrossfeed ()

