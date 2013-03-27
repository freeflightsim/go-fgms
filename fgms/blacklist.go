package fgms

import (
	"log"
	"net"
)


// --------------------------------------------------------
// AddBlacklist - Add an IP to the blacklist - (after DNS lookup)
func (me *FG_SERVER) AddBlacklist(FourDottedIP string) {
	log.Println("> Add Blacklist = ", FourDottedIP)
	
	// Do Checks in background
	go func(ip_str string){
		
		// Check DNS entry
		addrs, err := net.LookupHost(ip_str)
		if err != nil{
			log.Println("\tFAIL: Blacklist - No IP address for address = ", ip_str)
			return 
		}
		log.Println("\tOK:   Blacklist Added -  DNS Lookup OK: ", ip_str, addrs, ip_str == addrs[0])
		me.BlackList[ addrs[0] ] = true
	}(FourDottedIP)
} 



// Check if the client is black listed. true if blacklisted
func (me *FG_SERVER) IsBlackListed(SenderAddress *NetAddress) bool {
	_, found :=  me.BlackList[SenderAddress.IpAddress]
	if found {
		return true
	}
	return false
} 


