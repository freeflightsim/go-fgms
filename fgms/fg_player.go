package fgms

import (
	"net"


	"github.com/freeflightsim/go-fgms/message"
)

type FG_Player struct {
	
	//Origin string
	Address *net.UDPAddr  `json:"-"`
	//Conn *net.UDPConn
	
	Callsign string `json:"callsign"`
	//Passwd string
	ModelName string
	Aircraft string `json:"model"`
	
	JoinTime int64
	Timestamp int64
	
	LastPos Point3D `json:"-"`
	LastOrientation Point3D `json:"-"`
	
	IsLocal bool
	
	Error string `json:"-"`
	HasErrors bool `json:"-"`
	
	ClientID int
	LastRelayedToInactive int64  `json:"-"`
	
	// Packets recieved from client 
	PktsReceivedFrom uint  
	
	// Packets sent to client 
	PktsSentTo uint       
	
	// Packets from client sent to other players/relays 
	PktsForwarded uint    

} // FG_Player



// Returns the /model/747.400/AIRCRAFT.xml part
func (me *FG_Player) UpdatePosition(position *message.PositionMsg)  {

}


// Creates a new FG_Player object with timestamp set
func DEADNewFG_Player() *FG_Player {
	ob := new(FG_Player)
	ob.Timestamp = Now()
	ob.JoinTime  = ob.Timestamp
	return ob
}

