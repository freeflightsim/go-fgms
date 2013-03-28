

package fgms

import (
	"fmt"
	"net"
)

type FG_Conn struct{
	HostName string
	Port int
	Url string
	Ip string
	Active bool
	LastError string
	Conn *net.UDPConn
}



func NewFG_Conn(host_name string, port int) *FG_Conn {
	ob := new(FG_Conn)
	ob.HostName = host_name
	ob.Port = port
	ob.Url = fmt.Sprintf("%s:%d", host_name, port)
	return ob
}