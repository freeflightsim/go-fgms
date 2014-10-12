
package fgms

import(
	"fmt"
	"log"
	"net"
	"time"
)

import(
	"github.com/davecgh/go-xdr/xdr"
	"github.com/FreeFlightSim/go-fgms/flightgear"
)



//  Adds a new crossfeed server into internal list - after resolution of address etc
func (me *FG_SERVER) AddCrossfeed( host_name string, port int){

	log.Println("> Add Crossfeed = ", host_name, port)		
	
	// Create object
	conn := NewFG_Conn(host_name, port)
	me.Crossfeeds[conn.Url] = conn
	
	// Go check and setup
	go me.InitSetupCrossfeed(conn)
	
} 


//= Starts a timer to check servers that are down ie Active = false (currently every 60 secs)
// this is started in Loop() as a goroutine go me.StartCrossfeedCheckTimer()
func (me *FG_SERVER) StartCrossfeedCheckTimer(){
	
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

// Attempt to connect and setup a Crossfeed Connection
// Should dns fail, address not exist or not able to connect
// the connection witll be marked as conn.Active = false with conn.LastError
func (me *FG_SERVER) InitSetupCrossfeed( conn *FG_Conn){

	log.Println("> InitSetupCrossfeed = ", conn.Url )
		
	if conn.Active {
		return // we need to define when its inactive, eg connection dropped
	}		
		
	//= resolve with UDP address			
	udp_addr, err := net.ResolveUDPAddr("udp", conn.Url)
	if err != nil {
		conn.Active = false
		conn.LastError = "Could nto resolve IP address"
		log.Println("\tFAIL: Crossfeed to resolve UDP address:  ", conn.Url, err)
		return
	}
		
	//= open socket and listen
	var err_listen error
	conn.Conn, err_listen = net.ListenUDP("udp", udp_addr)
	if err_listen != nil {
		conn.Active = false
		conn.LastError = "Couldnt open UDP port"
		log.Println("\tFAIL: Crossfeed FAIL to Open:  ", conn.Url, udp_addr, err_listen)
		return
	}
	conn.Active = true
	conn.LastError = ""
	log.Println("\tOK:   Crossfeed Added -  ", conn.Url, udp_addr, err_listen)
	
} // FG_SERVER::AddCrossfeed()



/*   Send message to all crossfeed servers.
	Crossfeed servers receive all traffic without condition,
	mainly used for testing and debugging, and crossfeed.fgx.ch
	http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_server.cxx#line1154
*/
func (me *FG_SERVER) SendToCrossfeed(Msg []byte, Bytes int, SenderAddress *net.UDPAddr){

	//T_MsgHdr*       MsgHdr;
	//uint32_t        MsgMagic;
	//int             sent;
	
	//MsgHdr    = (T_MsgHdr *) Msg;
	//MsgMagic  = MsgHdr->Magic;
	//MsgHdr->Magic = XDR_encode<uint32_t> (RELAY_MAGIC);
	
	// Not sure what is happening, but we create another payload, by unmarshalling Msg again ?
	var MsgHdr flightgear.T_MsgHdr
	_, err := xdr.Unmarshal(Msg, &MsgHdr)
	if err != nil {
		fmt.Println("XDR Decode Error in SendToCrossfeed - Should never happen?", err)
		return
	}
	MsgHdr.Magic = RELAY_MAGIC
	
	encoded, err := xdr.Marshal(MsgHdr)
	if err != nil {
		fmt.Println("XDR Encode Error in SendToCrossfeed - Should never happen?", err)
		return
	}
	//mT_RelayListIt CurrentCrossfeed = m_CrossfeedList.begin();
	//while (CurrentCrossfeed != m_CrossfeedList.end())
	//{
	for _, loopCF := range me.Crossfeeds {

		//if (CurrentCrossfeed->Address.getIP() != SenderAddress.getIP()) ?? But same address different port ??
		_, err := loopCF.Conn.WriteToUDP(encoded, SenderAddress)
		if err != nil {
			me.CrossFeedFailed++
		}else {
			me.CrossFeedSent++
		}
		//CurrentCrossfeed++;
	}
	//MsgHdr->Magic = MsgMagic;  // restore the magic value ? umm not used now ?
} // FG_SERVER::SendToCrossfeed ()

