
package fgms

import(
	"net"
	"strings"
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


//////////////////////////////////////////////////
// mT_Relay - Type of list of relays
type MT_Relay struct {
	Name string
	Address *NetAddress // TODO = netAddress  Address
}

func NewMT_Relay(host_name string, port int) *MT_Relay {
	ob := new(MT_Relay)
	ob.Name = strings.Split(host_name, ".")[0]
	ob.Address = &NetAddress{Host: host_name, Port: port}
	return ob
}