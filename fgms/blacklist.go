package fgms

import (
	"log"
	"net"
)


// AddBlAdds an IP to the blacklist after a succesful go DNS lookup
func (me *FgServer) AddBlacklist(host_name string) {
	
	log.Println("> Add Blacklist = ", host_name)
	
	// Do Checks in background
	go func(host_name string){
		
		// Check DNS entry
		addrs, err := net.LookupHost(host_name)
		if err != nil{
			log.Println("\tFAIL: Blacklist - No IP address for address = ", host_name)
			return 
		}
		log.Println("\tOK:   Blacklist Added - DNS Lookup: ", host_name, " = ", addrs[0], host_name == addrs[0])
		me.BlackList[ addrs[0] ] = true
	}(host_name)
} 



// Check if the client is black listed. true if blacklisted TODO Fix ME
func (me *FgServer) IsBlackListed(SenderAddress *NetAddress) bool {
	_, found :=  me.BlackList[SenderAddress.IpAddress]
	if found {
		return true
	}
	return false
} 


