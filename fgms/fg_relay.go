
package fgms


type NetAddress struct {
	host string
	port int
	ip	string
	family int
	broadcase bool
}

func (me *NetAddress) GetIP() string {

	return "IPACCREDD"
}


//////////////////////////////////////////////////
// mT_Relay - Type of list of relays
type MT_Relay struct {
	Name string
	Address NetAddress // TODO = netAddress  Address
}

func NewMT_Relay(host_name string, port int) *MT_Relay {
	ob := new(MT_Relay)
	ob.Name = host_name
	
	return ob
}