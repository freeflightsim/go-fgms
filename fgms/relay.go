
package fgms

import(
	"fmt"
	"log"
	"net"
	//"strings"
)

type RelayData struct {
	Bytes []byte
	Player *Player
}

type relays struct {
	Chan chan RelayData
	Hosts map[string]net.Conn
	PktsForwarded int64
}

var Relays *relays

func SetupRelays() {

	Relays = new(relays)
	Relays.Chan = make(chan RelayData)
	Relays.Hosts = make(map[string]net.Conn, 0)
	go Relays.Listen()
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
		udp_addr, err := net.ResolveUDPAddr("udp4", host_port)
		if err != nil {
			log.Println("\tFAIL: Relay - failed to resolve UDP address  ", host_port, udp_addr, err)
			return
		}

		//= Now we open socket and listen
		var err_listen error
		me.Hosts[host_port], err_listen = net.Dial("udp4", udp_addr.String())
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





// Listen on channel and send to relays
func (me *relays) Listen(){

	for {
		// Got data from channel
		relay_data := <- me.Chan

		now := Now()
		UpdateInactive := (now - relay_data.Player.LastRelayedToInactive) > UPDATE_INACTIVE_PERIOD
		if UpdateInactive {
			relay_data.Player.LastRelayedToInactive = now
		}

		for _, host := range me.Hosts {
			if UpdateInactive { //|| IsInRange(*relay, *SendingPlayer) {
				fmt.Println("relay to=",  host)
				//if (CurrentRelay->Address.getIP() != SendingPlayer->Address.getIP())
				//{
				//  m_DataSocket->sendto(Msg, Bytes, 0, &CurrentRelay->Address);
				//host.WriteToUDP(relay_data.Bytes, relay_data.Client.Address)
				host.Write(relay_data.Bytes)
				me.PktsForwarded++
				// }
				//}
				//CurrentRelay++;
			}
		}
	}
}
