

package fgms

import (
	"fmt"
	"net"
)

type UDP_Conn struct{
	HostName string
	Port int
	Url string
	Ip string
	Active bool
	LastError string
	Sock net.Conn
}



func NewUDPConn(host_name string, port int) *UDP_Conn {
	ob := new(UDP_Conn)
	ob.HostName = host_name
	ob.Port = port
	ob.Url = fmt.Sprintf("%s:%d", host_name, port)
	return ob
}
