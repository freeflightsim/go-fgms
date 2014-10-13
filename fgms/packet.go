package fgms


import(
	"fmt"
	"net"
	"strings"

	"github.com/freeflightsim/go-fgms/message"
)


//------------------------------------------------------------------------

// Handle client connections
func (me *FgServer) HandlePacket(xdr_bytes []byte, length int, sender_address *net.UDPAddr){
	
	//T_MsgHdr*       MsgHdr;
	//var MsgHdr message.T_MsgHdr
	//T_PositionMsg*  PosMsg;
	//var PosMsg flightgear.T_PositionMsg
	
	//uint32_t        MsgId;
	//uint32_t        MsgMagic;
	//Timestamp time.Time

	
	//var SenderPosition Point3D
	//var SenderOrientation Point3D
	//Point3D         PlayerPosGeod;
	//mT_PlayerListIt CurrentPlayer;
	//mT_PlayerListIt SendingPlayer;
	//unsigned int    PktsForwarded = 0;
	
	//Timestamp = time.Now() //time(0);
	//MsgHdr    = (T_MsgHdr *) Msg;
	//MsgHdr :=  
	
	//fmt.Println("MSG=", len(Msg))
	//var header message.HeaderMsg

	if Blacklist.IsBlackListed(sender_address){
		//me.BlackListRejected++
		fmt.Println("Blacklisted")
		return
	}

	// Decode header message, exit if error
	header, remainingBytes, err := message.DecodeHeader(xdr_bytes)
	if err != nil{
		fmt.Println("XDR header error", err)
		me.PacketsInvalid++
		return
	}
	//fmt.Println("remain=", len(remainingBytes), address.String(), header.Callsign())
	me.PacketsReceived++


	//timestamp := Now()

	CrossFeed.Chan <- xdr_bytes

	//me.SendToCrossfeed(xdr_bytes, sender_address)
	//Crossfeeds.Chan <- xdr_bytes
	//------------------------------------------------------
	// First of all, send packet to all crossfeed servers.
	//SendToCrossfeed (Msg, Bytes, SenderAddress); ?? SHould then be send pre vaildation ?
	//me.SendToCrossfeed(Msg, Bytes, SenderAddress)



	//------------------------------------------------------
	//=  Now do the local processing TODO
	//if me.IsBlackListed(SenderAddress) {
	//	me.BlackListRejected++
	//	return
	//}

	if header.Magic == message.RELAY_MAGIC { // not a local client
		if !Relays.IsKnown(sender_address) {
			me.UnknownRelay++ 
			return
		}else{
			me.RelayMagic++ // bump relay magic packet
		}
	}

	callsign := header.Callsign()

	var player *FG_Player
	var position message.PositionMsg
	var exists bool
	var remBytes []byte
	var err_pos error

	// Check if entry exists
	player, exists = me.Players[callsign]

	if exists == false &&  header.Type  != message.TYPE_POS {
		// ignore until a position message
		return
	}
	if exists == true && player.Address.String() != sender_address.String() {
		// sender has same callsign but different address, so ignore
		return
	}


	// Decode position packer
	if header.Type == message.TYPE_POS	{

		position, remBytes, err_pos = message.DecodePosition(remainingBytes)
		if err != nil{
			fmt.Println("XDR Decode Position Error", err_pos)
			return
		}else if 1 == 2 {
			fmt.Println("remain2=", len(remBytes))
		}

		if position.Position[X] == 0.0 || position.Position[Y] == 0.0 || position.Position[Z] == 0.0 {
			return // ignore while position is not settled
		}
		me.PositionData++
	} else {
		me.NotPosData++
	}

	// Create new player
	if exists == false {
		player = me.AddClient(&header, &position, sender_address)
	}

	fmt.Println( callsign, position.Position[X], position.Position[Y])
	player.UpdatePosition(&position)
	player.Timestamp = Now()
	player.PktsReceivedFrom++

	//////////////////////////////////////////
	//
	//      send the packet to all clients.
	//      since we are walking through the list,
	//      we look for the sending client, too. if it
	//      is not already there, add it to the list
	//
	//////////////////////////////////////////////////

	isObserver :=  strings.ToLower(callsign)[0:3] ==  "obs"
	for _, lp := range me.Players {
		
		//= ignore clients with errors
		if lp.HasErrors {
			continue // Umm is this locked out forever ?
		}
		
		
		// Sender == CurrentPlayer?
		/*   FIXME: if Sender is a Relay,
					CurrentPlayer->Address will be
				address of Relay and not the client's!
				so use a clientID instead
		*/
		//if loopCallsign == callsign { // alterative == CurrentPlayer.Callsign == xCallsign
			//if header.Type == message.TYPE_POS	{
				// Update this players position
				//player.LastPos.Set( position.Position[X], position.Position[Y], position.Position[Z])
				//player.LastOrientation.Set( float64(position.Orientation[X]), float64(position.Orientation[Y]), float64(position.Orientation[Z]))
			//	loopPlayer.LastPos.Set( position.Position[X], position.Position[Y], position.Position[Z])
			//	loopPlayer.LastOrientation.Set( float64(position.Orientation[X]), float64(position.Orientation[Y]), float64(position.Orientation[Z]))
			//}//else{
				//SenderPosition    = loopPlayer.LastPos
				//SenderOrientation = loopPlayer.LastOrientation
			//}
			//SendingPlayer = CurrentPlayer
			//loopPlayer.Timestamp = timestamp
			//loopPlayer.PktsReceivedFrom++
			//CurrentPlayer++;
			//continue; // don't send packet back to sender
		//}
		///     do not send packets to clients if the
		//      origin is an observer, but do send
		//      chat messages anyway
		//      FIXME: MAGIC = SFGF!
		if isObserver && header.Type != message.TYPE_CHAT {
			continue
		}
		
		// Do not send packet to clients which  are out of reach.
		//if xIsObserver == false && int(Distance(SenderPosition, loopPlayer.LastPos)) > me.PlayerIsOutOfReach {
			//if ((Distance (SenderPosition, CurrentPlayer->LastPos) > m_PlayerIsOutOfReach)
			//&&  (CurrentPlayer->Callsign.compare (0, 3, "obs", 3) != 0))
			//{
			//CurrentPlayer++ 
			//continue
		//}
		
		//  only send packet to local clients
		if lp.IsLocal && lp != player {
			//SendChatMessages (CurrentPlayer);
			//m_DataSocket->sendto (Msg, Bytes, 0, &CurrentPlayer->Address);
			_, err := me.DataSocket.WriteToUDP(xdr_bytes, player.Address)
			if err != nil {
				// TODO ?
			}
			lp.PktsSentTo++
			me.PktsForwarded++
		}
		//CurrentPlayer++; 
		//
	} 
	/* 
	if (SendingPlayer == m_PlayerList.end())
	{ // player not yet in our list
		// should not happen, but test just in case
		SG_LOG (SG_SYSTEMS, SG_ALERT, "## BAD => "
		<< MsgHdr->Callsign << ":" << SenderAddress.getHost()
		<< " : " << SenderIsKnown (MsgHdr->Callsign, SenderAddress)
		);
		return;
	}
	DeleteMessageQueue ();
	*/
	//SendingPlayer := NewFG_Player() // placleholder
	//me.SendToRelays (xdr_bytes, length, player)
	Relays.Chan <- RelayData{Bytes: xdr_bytes, Client: player}
	
} // FgServer::HandlePacket ( char* sMsg[MAX_PACKET_SIZE] )



