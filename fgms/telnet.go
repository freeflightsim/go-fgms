

package fgms

import(
	"net"
)
 
type TelnetServer  struct {
	Port int
	Reinit bool
	Received int
	Conn *net.Conn
}

func NewTelnetServer() *TelnetServer {
	ob := new(TelnetServer)
	return ob
}