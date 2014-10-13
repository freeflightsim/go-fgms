package fgms

import (
	"log"
	"net"
)

type blacklist struct {
	Hosts map[string]bool
	Rejected int64
}

var Blacklist blacklist


func SetupBlackList(){
	//log.Println("InInitBlackList>>>>>>>>>>>>>>>")
	Blacklist = blacklist{}
	Blacklist.Hosts = make(map[string]bool, 0)
	//log.Println("InitBlacklist")
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
func (me *blacklist) IsBlackListed(address *net.UDPAddr) bool {

	host, _, err := net.SplitHostPort(address.String())
	if err != nil {
		return false // ??
	}
	_, found := Blacklist.Hosts[host]
	if found {
		me.Rejected++
	}
	return found

}
