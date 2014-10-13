package fgms

import (
	"net"

	"github.com/freeflightsim/go-fgms/message"
)

type Player struct {
	
	//Origin string
	Address *net.UDPAddr  `json:"-"`
	//Conn *net.UDPConn
	
	Callsign string `json:"callsign"`
	//Passwd string
	ModelName string `json:"model_name"`
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



// Update Current Pos
func (me *Player) UpdatePosition(position *message.PositionMsg)  {
	me.LastPos.Set( position.Position[X], position.Position[Y], position.Position[Z])
	me.LastOrientation.Set( float64(position.Orientation[X]), float64(position.Orientation[Y]), float64(position.Orientation[Z]))
}

// return Geod positions
func (me *Player) LatLonAlt() (float64, float64, float64)  {
	xp := SG_CartToGeod(me.LastPos)
	return xp.X, xp.Y, xp.Z

}

