

package fgms

import(
	"net"
)
 
// TelnetServer container
type TelnetServer  struct {
	Port int
	Reinit bool
	Received int
	Conn *net.Conn
}

// Constructs and return TelnetServer
func NewTelnetServer() *TelnetServer {
	ob := new(TelnetServer)
	return ob
}