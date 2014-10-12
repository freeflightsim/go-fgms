
package fgms

import(
	"fmt"
	"log"
	"net"
	"time"
)



// Insert a new relay server into internal list (does a DNS lookup)
func (me *FgServer) AddRelay(host_name string, port int) {
	
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
	
} // FgServer::AddRelay()




//////////////////////////////////////////////////////////////////////
//  Check if the sender is a known relay, return true if known relay
func (me *FgServer) IsKnownRelay(senderAddress *net.UDPAddr) bool{
	fmt.Println("IsKnownRelay", senderAddress.String())
	/*mT_RelayListIt  CurrentRelay = m_RelayList.begin();
	while (CurrentRelay != m_RelayList.end())
	{
		if (CurrentRelay->Address.getIP() == SenderAddress.getIP())
		{
		return (true);
		}
		CurrentRelay++;
	}*/
	//_, ok := me.RelayMap[senderAddress.String()]
	//if ok {
	//	return true
	//}

	//string ErrorMsg;
	//ErrorMsg  = SenderAddress.getHost();
	//ErrorMsg += " is not a valid relay!";
	//me.AddBlacklist(senderAddress.IpAddress)
	//SG_LOG (SG_SYSTEMS, SG_ALERT, "UNKNOWN RELAY: " << ErrorMsg);
	return false
} // FgServer::IsKnownRelay ()




// ---------------------------------------------------------

// Send message to all relay servers

func (me *FgServer) SendToRelays(Msg []byte, Bytes int , SendingPlayer *FG_Player){

	//T_MsgHdr*       MsgHdr;
	//uint32_t        MsgMagic;
	var PktsForwarded uint = 0
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
	for idx, relayConn := range me.Relays {
		if UpdateInactive { //|| IsInRange(*relay, *SendingPlayer) {
			fmt.Println("relay to=", idx, relayConn)
			//if (CurrentRelay->Address.getIP() != SendingPlayer->Address.getIP())
			//{
			//  m_DataSocket->sendto(Msg, Bytes, 0, &CurrentRelay->Address);
			relayConn.WriteToUDP(Msg, SendingPlayer.Address)
			me.PktsForwarded++
			// }
			//}
			//CurrentRelay++;
		}
	}
	SendingPlayer.PktsForwarded += PktsForwarded
	//MsgHdr->Magic = XDR_encode<uint32_t> (MsgMagic);  // restore the magic value
} // FgServer::SendToRelays ()



