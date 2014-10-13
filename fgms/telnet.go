

package fgms

import(
	"bytes"
	"net"
	"log"
	"strconv"
)

 import(
	//"github.com/fgx/go-fgms/flightgear"
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



//---------------------------------------------------------------------------

//  Handle a telnet session. if a telnet connection is opened, this 
// method outputs a list  of all known clients.
func (me *FgServer) HandleTelnetData(conn net.Conn){

	//var errno int = 0
	var Message string  = ""
	//buf := make([]byte, 4096)

	
	/** @brief  Geodetic Coordinates */
	//Point3D         PlayerPosGeod;  
	//FG_Player CurrentPlayer;
	//netSocket       NewTelnet;
	//unsigned int  it;
	//NewTelnet.setHandle (Fd);
	//errno = 0;
	//////////////////////////////////////////////////
	//
	//      create the output message
	//      header
	//
	//////////////////////////////////////////////////
	Message  = "# This is " + me.ServerName 
	Message += "\n"
	Message += "# FlightGear Multiplayer Server version: " + VERSION
	Message += "\n"
	Message += "# using protocol version: "
	Message += GetProtocolVersionString() // FIX ME PLEASE
	Message += " (LazyRelay enabled)"
	Message += "\n"
	//buf.Add

	// print conn.RemoteAddr()
	/* if ( m_IsTracked )
	{
		Message += "# This server is tracked: ";
		Message += m_Tracker->GetTrackerServer();
		Message += "\n";
	}
	if (NewTelnet.send (Message.c_str(),Message.size(), MSG_NOSIGNAL) < 0)
	{
		if ((errno != EAGAIN) && (errno != EPIPE))
		{
		SG_LOG (SG_SYSTEMS, SG_ALERT, "FgServer::HandleTelnet() - " << strerror (errno));
		}
		return (0);
	} */
	/* pthread_mutex_lock (& m_PlayerMutex); 
	Message  = "# "+ NumToStr (m_PlayerList.size(), 0);
	pthread_mutex_unlock (& m_PlayerMutex);
	Message += " pilot(s) online\n";
	if (NewTelnet.send (Message.c_str(),Message.size(), MSG_NOSIGNAL) < 0)
	{
		if ((errno != EAGAIN) && (errno != EPIPE))
		{
		SG_LOG (SG_SYSTEMS, SG_ALERT, "FgServer::HandleTelnet() - " << strerror (errno));
		}
		return (0);
	}*/
	Message += " pilot(s) online\n"
	

	//== Create list of players
	for callsign, CurrentPlayer := range me.Players {
	//it = 0;
	//for (;;)
	//{
		//pthread_mutex_lock (& m_PlayerMutex);
		//if (it < m_PlayerList.size())
		//{
		//CurrentPlayer = m_PlayerList[it]; 
		//it++;
		//}
		//else
		//{
		//pthread_mutex_unlock (& m_PlayerMutex);
		//break;
		//}
		//pthread_mutex_unlock (& m_PlayerMutex);
		//TODO sgCartToGeod (CurrentPlayer.LastPos, PlayerPosGeod);
		
		line  := callsign + "@"
		//Message += CurrentPlayer.Callsign + "@"
		if CurrentPlayer.IsLocal {
			line += "LOCAL: "
		}else{
			//mT_RelayMapIt Relay = m_RelayMap.find(CurrentPlayer.Address.getIP())
			//if (Relay != m_RelayMap.end()){
			//	line += Relay->second + ": "
			//}else{
			//	line += CurrentPlayer.Origin + ": "
			//}
		}
		if CurrentPlayer.Error != "" {
			line += CurrentPlayer.Error + " "
		}
		PlayerPosGeod := SG_CartToGeod (CurrentPlayer.LastPos)
		//= Last Position
		//Message += NumToStr (CurrentPlayer.LastPos[X], 6)+" ";
		//Message += NumToStr (CurrentPlayer.LastPos[Y], 6)+" ";
		//Message += NumToStr (CurrentPlayer.LastPos[Z], 6)+" ";
		//
		// http://golang.org/pkg/strconv/#FormatFloat
		Message += strconv.FormatFloat( CurrentPlayer.LastPos.X, 'f', 6, 32)  + " "
		Message += strconv.FormatFloat( CurrentPlayer.LastPos.Y, 'f', 6, 32)  + " "
		Message += strconv.FormatFloat( CurrentPlayer.LastPos.Z, 'f', 6, 32)  + " "
		
		//= Lat/Lon/Alt
		//Message += NumToStr (PlayerPosGeod[Lat], 6)+" ";
		//Message += NumToStr (PlayerPosGeod[Lon], 6)+" ";
		//Message += NumToStr (PlayerPosGeod[Alt], 6)+" ";
		Message += strconv.FormatFloat( PlayerPosGeod.Lat(), 'f', 6, 32)  + " "
		Message += strconv.FormatFloat( PlayerPosGeod.Lon(), 'f', 6, 32)  + " "
		Message += strconv.FormatFloat( PlayerPosGeod.Alt(), 'f', 6, 32)  + " "
		
		//Message += NumToStr (CurrentPlayer.LastOrientation[X], 6)+" ";
		//Message += NumToStr (CurrentPlayer.LastOrientation[Y], 6)+" ";
		//Message += NumToStr (CurrentPlayer.LastOrientation[Z], 6)+" ";
		
		Message += strconv.FormatFloat( CurrentPlayer.LastOrientation.X, 'f', 6, 32)  + " "
		Message += strconv.FormatFloat( CurrentPlayer.LastOrientation.Y, 'f', 6, 32)  + " "
		Message += strconv.FormatFloat( CurrentPlayer.LastOrientation.Z, 'f', 6, 32)  + " "
		
		//Message += CurrentPlayer.ModelName;
		Message += CurrentPlayer.ModelName
		
		Message += "\n"
		Message += line
		
		/*if (NewTelnet.send (Message.c_str(),Message.size(), MSG_NOSIGNAL) < 0)
		{
			if ((errno != EAGAIN) && (errno != EPIPE))
			{
				SG_LOG (SG_SYSTEMS, SG_ALERT, "FgServer::HandleTelnet() - " << strerror (errno));
			}
			return (0);
		}*/
	}
	// NewTelnet.close ();
	var buffer bytes.Buffer
	buffer.WriteString( Message )
	_, err := conn.Write( buffer.Bytes() )
	if err != nil {
		log.Println("error", err)
	}
	conn.Close()
	//return (0);
} // FgServer::HandleTelnet ()

