package fgms

import (
	"log"
	"net"
)

type blacklist struct {
	Hosts map[string]bool // we could use slice, but usign map instead for ocnvenience
}

var Blacklist blacklist


func init(){

	Blacklist = blacklist{}
	Blacklist.Hosts = make(map[string]bool, 0)

}

// Add an IP to the blacklist after a succesful  DNS lookup
func (me *blacklist) Add(host_name string) {

	// Do checks in thread
	go func(host_name string){

		addrs, err := net.LookupHost(host_name)
		if err != nil{
			log.Println("Blacklist: Failed - No IP address for address = ", host_name)
			return
		}
		log.Println("Blacklist: Added ", host_name, " = ", addrs[0], host_name == addrs[0])
		me.Hosts[addrs[0]] = true
	}(host_name)
}

// Check if the client is black listed. true if blacklisted
func (me *blacklist) Contains(address *net.UDPAddr) bool {

	host, _, err := net.SplitHostPort(address.String())
	if err != nil {
		return false // ??
	}
	_, found := Blacklist.Hosts[host]
	return found

}


// AddBlAdds an IP to the blacklist after a succesful go DNS lookup
func (me *FgServer) AddBlacklist(host_name string) {
	
	//log.Println("> Add Blacklist = ", host_name)
	
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


