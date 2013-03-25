

package fgms

import(
	"net"
)

 
// TelnetServer container
type TelnetServer  struct {
	Addr string
	Port int
	Reinit bool
	Received int
	//Conn *net.Conn
	Listen net.Listener
}

// Constructs and return TelnetServer
func NewTelnetServer() *TelnetServer {
	ob := new(TelnetServer)
	return ob
}
