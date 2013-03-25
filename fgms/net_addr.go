
package fgms

import(
	"net"
)

type NetAddress struct {
	Host string
	Port int
	IpAddress	string
	Family int
	Broadcast bool
}

func (me *NetAddress) LookupIP() error {
	addrs, err := net.LookupHost(me.Host)
	if err != nil {
		return err
	}
	me.IpAddress = addrs[0]
	return nil
}
