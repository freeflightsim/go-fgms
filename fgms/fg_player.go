package fgms

import (
	"time"
)

type FG_Player struct {
	//public:
	Origin string
	Address *netAddress    
	Callsign string
	Passwd string
	ModelName string
	JoinTime time.Date
	Timestamp time.Date
	//Point3D       LastPos;
	//Point3D       LastOrientation;
	IsLocal bool
	Error string //;    // in case of errors
	HasErrors bool
	ClientID int
	LastRelayedToInactive time.Date
	
	// Packets recieved from client 
	PktsReceivedFrom uint  
	
	// Packets sent to client 
	PktsSentTo uint       
	
	//  Packets from client sent to other players/relays 
	PktsForwarded uint    

	//FG_Player ();
	//FG_Player ( const FG_Player& P);
	//~FG_Player ();
	//void operator =  ( const FG_Player& P );
	//private:
	//void assign ( const FG_Player& P );
} // FG_Player


