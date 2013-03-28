package fgms

import (
	"net"
	"path/filepath"
)

type FG_Player struct {
	
	//Origin string
	Address *net.UDPAddr    
	//Conn *net.UDPConn
	
	Callsign string // But this is also key so maybe unneeded ?
	Passwd string
	ModelName string
	
	JoinTime int64 // epoch
	Timestamp int64 // epoch
	
	LastPos Point3D
	LastOrientation Point3D
	
	IsLocal bool
	
	Error string //;    // in case of errors
	HasErrors bool
	
	ClientID int
	LastRelayedToInactive int64 
	
	// Packets recieved from client 
	PktsReceivedFrom uint  
	
	// Packets sent to client 
	PktsSentTo uint       
	
	// Packets from client sent to other players/relays 
	PktsForwarded uint    

} // FG_Player



// Returns the /model/747.400/AIRCRAFT.xml part
func (me *FG_Player) Aircraft() string {
	s := me.ModelName
	return filepath.Base(s) // ? TO CHECK  
}


// Creates a new FG_Player object with timestamp set
func NewFG_Player() *FG_Player {
	ob := new(FG_Player)
	ob.Timestamp = Now()
	ob.JoinTime  = ob.Timestamp
	return ob
}

