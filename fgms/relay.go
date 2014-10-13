
package fgms

import(
	"fmt"
	"log"
	"net"
	//"strings"
)

type relays struct {
	Chan chan []byte
	Hosts map[string]*net.UDPConn
	PktsForwarded int64
}

var Relays *relays

func SetupRelays() {

	Relays = new(relays)
	Relays.Chan = make(chan []byte)
	Relays.Hosts = make(map[string]*net.UDPConn, 0)
}


// Insert a new relay server into internal list (does a DNS lookup)
func (me *relays) Add(host_name string, port int) {

	log.Println("> Add Relay = ", host_name, port)


	//= Now go and do check is background
	go func(host_name string, port int){

		// Get IP address from DNS
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
		me.Hosts[host_port], err_listen = net.ListenUDP("udp", udp_addr)
		if err_listen != nil {
			log.Println("\tFAIL: Relay - Cannot listen on UDP  ", host_port, udp_addr, err_listen)
			return
		}
		log.Println("    < Relay Added OK  ", host_port, udp_addr, err_listen)
	}(host_name, port)

}

// Check if the address is a known relay
func (me *relays) IsKnown(address *net.UDPAddr) bool{
	_, found := me.Hosts[address.String()]
	return found
}



// ---------------------------------------------------------

// Send message to all relay servers

func (me *relays) Send(Msg []byte, Bytes int , sending_player *FG_Player){

	//T_MsgHdr*       MsgHdr;
	//uint32_t        MsgMagic;
	var PktsForwarded uint = 0
	//mT_RelayListIt  CurrentRelay;
	//time_t          Now;

	if !sending_player.IsLocal && !Server.IamHUB {
		return
	}
	//Now   = time (0);
	now := Now()
	//MsgHdr    = (T_MsgHdr *) Msg;
	//MsgMagic  = XDR_decode<uint32_t> (MsgHdr->Magic);
	//MsgHdr->Magic = XDR_encode<uint32_t> (RELAY_MAGIC);
	UpdateInactive := (now - sending_player.LastRelayedToInactive) > UPDATE_INACTIVE_PERIOD
	if UpdateInactive {
		sending_player.LastRelayedToInactive = now
	}
	//CurrentRelay = m_RelayList.begin();
	//while (CurrentRelay != m_RelayList.end())
	for idx, host := range me.Hosts {
		if UpdateInactive { //|| IsInRange(*relay, *SendingPlayer) {
			fmt.Println("relay to=", idx, host)
			//if (CurrentRelay->Address.getIP() != SendingPlayer->Address.getIP())
			//{
			//  m_DataSocket->sendto(Msg, Bytes, 0, &CurrentRelay->Address);
			host.WriteToUDP(Msg, sending_player.Address)
			me.PktsForwarded++
			// }
			//}
			//CurrentRelay++;
		}
	}
	sending_player.PktsForwarded += PktsForwarded
	//MsgHdr->Magic = XDR_encode<uint32_t> (MsgMagic);  // restore the magic value
} // FgServer::SendToRelays ()

// Listen for xdr packets from channel, and send to xrossfeeds
func (me *relays) Listen(){
	fmt.Println("Relays: Listening")
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
