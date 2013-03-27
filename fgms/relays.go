
package fgms

import(
	"fmt"
	"log"
	"net"
	"time"
)



// Insert a new relay server into internal list (does a DNS lookup)
func (me *FG_SERVER) AddRelay(host_name string, port int) {
	
	log.Println("> Add Relay = ", host_name, port)
	
	//= Now go and do check is background
	go func(host_name string, port int){
	
		//= Get IP address from DNS
		addrs, err := net.LookupHost(host_name)
		if err != nil {
			log.Println("\tFAIL: Relay - No IP address for Host ", host_name, addrs)
			return
		}
	
		//= Now resolve with UDP address			
		host_port := fmt.Sprintf("%s:%d", host_name, port)
		//log.Println("    < Relay - DNS Lookup OK:  ", host_name, addrs[0], s)
		udp_addr, err := net.ResolveUDPAddr("udp", host_port)
		if err != nil {
			log.Println("\tFAIL: Relay - failed to resolve UDP address  ", host_port, udp_addr, err)
			return
		}
		
		//= Now we open socket and listen
		var err_listen error
		me.Relays[host_port], err_listen = net.ListenUDP("udp", udp_addr)
		if err_listen != nil {
			log.Println("\tFAIL: Relay - Cannot listen on UDP  ", host_port, udp_addr, err_listen)
			return
		}
		log.Println("    < Relay Added OK  ", host_port, udp_addr, err_listen)
	}(host_name, port)	
	
} // FG_SERVER::AddRelay()







// ---------------------------------------------------------

// Send message to all relay servers

func (me *FG_SERVER) SendToRelays(Msg []byte, Bytes int , SendingPlayer *FG_Player){

//T_MsgHdr*       MsgHdr;
//uint32_t        MsgMagic;
//unsigned int    PktsForwarded = 0;
//mT_RelayListIt  CurrentRelay;
//time_t          Now;

if !SendingPlayer.IsLocal && !me.IamHUB {
	return
}
//Now   = time (0);
Now := time.Now().Unix()
//MsgHdr    = (T_MsgHdr *) Msg;
//MsgMagic  = XDR_decode<uint32_t> (MsgHdr->Magic);
//MsgHdr->Magic = XDR_encode<uint32_t> (RELAY_MAGIC);
UpdateInactive := (Now - SendingPlayer.LastRelayedToInactive) > UPDATE_INACTIVE_PERIOD
if UpdateInactive {
		SendingPlayer.LastRelayedToInactive = Now
}
//CurrentRelay = m_RelayList.begin();
//while (CurrentRelay != m_RelayList.end())
for idx, relay := range me.Relays {
	if UpdateInactive { //|| IsInRange(*relay, *SendingPlayer) {
		fmt.Println("relay to=", idx, relay)
		//if (CurrentRelay->Address.getIP() != SendingPlayer->Address.getIP())
		//{
		//  m_DataSocket->sendto(Msg, Bytes, 0, &CurrentRelay->Address);
		//  PktsForwarded++;
		// }
		//}
		//CurrentRelay++;
	}
}
//SendingPlayer->PktsForwarded += PktsForwarded;
//MsgHdr->Magic = XDR_encode<uint32_t> (MsgMagic);  // restore the magic value
} // FG_SERVER::SendToRelays ()



