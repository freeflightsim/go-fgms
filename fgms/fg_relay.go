
package fgms

import(
	
	"strings"
)



//////////////////////////////////////////////////
// mT_Relay - Type of list of relays
type MT_Relay struct {
	Name string
	Address NetAddress // TODO = netAddress  Address
}

func NewMT_Relay(host_name string, port int) *MT_Relay {
	ob := new(MT_Relay)
	ob.Name = strings.Split(host_name, ".")[0]
	ob.Address = NetAddress{Host: host_name, Port: port}
	return ob
}