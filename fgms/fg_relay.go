
package fgms

import(
	
	"strings"
)



//////////////////////////////////////////////////
// mT_Relay - Type of list of relays
/*type MT_Relay struct {
	Name string
	Address NetAddress // TODO = netAddress  Address
}
*/
func DEADNewMT_Relay(hostName string, port int) *NetAddress {
	ob := new(NetAddress)
	ob.Host = hostName
	ob.Port = port
	ob.NickName = strings.Split(hostName, ".")[0]
	
	return ob
}