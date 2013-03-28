
package fgms

import(
	//"net"
)
// Represent a Network address

type NetAddress struct {
	NickName string
	Host string
	Port int
	IpAddress	string
	Family int
	Broadcast bool
}
/*
func (me *NetAddress) LookupIP() error {
	addrs, err := net.LookupHost(me.Host)
	if err != nil {
		return err
	}
	me.IpAddress = addrs[0]
	return nil
}
*/