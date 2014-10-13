package fgms


import(
	//"bytes"
	"fmt"
	//"log"
	"net"		
	"strings"
	//"strconv"
	//"time"
	"unsafe"
)

import(
	//"github.com/davecgh/go-xdr/xdr"

	"github.com/FreeFlightSim/go-fgms/message"
	//"github.com/FreeFlightSim/go-fgms/flightgear"
)


//------------------------------------------------------------------------

// Handle client connections
func (me *FgServer) HandlePacket(xdr_bytes []byte, length int, address *net.UDPAddr){
	
	//T_MsgHdr*       MsgHdr;
	//var MsgHdr message.T_MsgHdr
	//T_PositionMsg*  PosMsg;
	//var PosMsg flightgear.T_PositionMsg
	
	//uint32_t        MsgId;
	//uint32_t        MsgMagic;
	//Timestamp time.Time

	
	var SenderPosition Point3D
	var SenderOrientation Point3D
	//Point3D         PlayerPosGeod;
	//mT_PlayerListIt CurrentPlayer;
	//mT_PlayerListIt SendingPlayer;
	//unsigned int    PktsForwarded = 0;
	
	//Timestamp = time.Now() //time(0);
	//MsgHdr    = (T_MsgHdr *) Msg;
	//MsgHdr :=  
	
	//fmt.Println("MSG=", len(Msg))
	//var header message.HeaderMsg


	// Decode header message, exit if error
	header, remainingBytes, err := message.DecodeHeader(xdr_bytes)
	if err != nil{
		fmt.Println("XDR header error", err)
		return
	}
	fmt.Println("remain=", len(remainingBytes), address.String())

	me.PacketsInvalid++

	timestamp := Now()


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
		if !me.IsKnownRelay(address) {
			me.UnknownRelay++ 
			return
		}else{
			me.RelayMagic++ // bump relay magic packet
		}
	}

	callsign := header.Callsign()

	var player *FG_Player
	var exists bool

	// Check if entry exists
	player, exists = me.Players[callsign]

	if exists == false &&  header.Type  != message.TYPE_POS {
		// ignore until a position message
		return
	}
	if exists == true && player.Address.String() != address.String() {
		// sender has same callsign but different address, do ignore
		return
	}

	if exists == false {
		// create new player
		player = new(FG_Player)
		player.Address = address
		player.Callsign = callsign
		player.IsLocal = header.Magic != message.RELAY_MAGIC
		me.Players[callsign] = player
	}




	//////////////////////////////////////////////////
	//    Store senders position
	//////////////////////////////////////////////////

	//var PosMsg message.PositionMsg
	if header.Type == message.TYPE_POS	{
		me.PositionData++
		position, remBytes, errPos := message.DecodePosition(remainingBytes)
		//remainingBytes2, errPos := xdr.Unmarshal(remainingBytes, &PosMsg)
		if err != nil{
			fmt.Println("XDR Decode Position Error", errPos)
			return
		}
		if 1 == 2 {
			fmt.Println("remain2=", len(remBytes), position.Model)
		}
		//PosMsg = (T_PositionMsg *) (Msg + sizeof(T_MsgHdr));
		//double x = XDR_decode64<double> (PosMsg->position[X]);
		//double y = XDR_decode64<double> (PosMsg->position[Y]);
		//double z = XDR_decode64<double> (PosMsg->position[Z]);
		x := position.Position[X]
		y := position.Position[Y]
		z := position.Position[Z]
		if x == 0.0 || y == 0.0 || z == 0.0 { // ignore while position is not settled
			return
		}
		//SenderPosition.Set (x, y, z);
		
		/* SenderOrientation.Set (
			XDR_decode<float> (PosMsg->orientation[X]),
			XDR_decode<float> (PosMsg->orientation[Y]),
			XDR_decode<float> (PosMsg->orientation[Z])
		)*/
		//TODO Wrong TYPE wtf!
		//SenderOrientation.Set(PosMsg.Orientations[X], PosMsg.Orientation[Y],	PosMsg.Orientation[Z])
		SenderOrientation.Set(0,0,0)
	} else {
		me.NotPosData++
	} 
	

	
	//////////////////////////////////////////
	//
	//      send the packet to all clients.
	//      since we are walking through the list,
	//      we look for the sending client, too. if it
	//      is not already there, add it to the list
	//
	//////////////////////////////////////////////////
	// MsgHdr->Magic = XDR_encode<uint32_t> (MSG_MAGIC);
	//SendingPlayer = m_PlayerList.end();
	//CurrentPlayer = m_PlayerList.begin();
	//while (CurrentPlayer != m_PlayerList.end())
	//{ 
	xCallsign := header.Callsign()
	xIsObserver :=  strings.ToLower(header.Callsign())[0:3] ==  "obs"
	for loopCallsign, loopPlayer := range me.Players {
		
		//= ignore clients with errors
		if loopPlayer.HasErrors {
			continue // Umm is this locked out forever ?
		}
		
		
		// Sender == CurrentPlayer?
		/*   FIXME: if Sender is a Relay,
					CurrentPlayer->Address will be
				address of Relay and not the client's!
				so use a clientID instead
		*/
		if loopCallsign == xCallsign { // alterative == CurrentPlayer.Callsign == xCallsign 
			if header.Type == message.TYPE_POS	{
				loopPlayer.LastPos         = SenderPosition
				loopPlayer.LastOrientation = SenderOrientation
			}else{
				SenderPosition    = loopPlayer.LastPos
				SenderOrientation = loopPlayer.LastOrientation
			}
			//SendingPlayer = CurrentPlayer
			loopPlayer.Timestamp = timestamp
			loopPlayer.PktsReceivedFrom++
			//CurrentPlayer++;
			continue; // don't send packet back to sender
		}
		///     do not send packets to clients if the
		//      origin is an observer, but do send
		//      chat messages anyway
		//      FIXME: MAGIC = SFGF!
		if xIsObserver && header.Type != message.TYPE_CHAT {
			return
		}
		
		// Do not send packet to clients which  are out of reach.
		if xIsObserver == false && int(Distance(SenderPosition, loopPlayer.LastPos)) > me.PlayerIsOutOfReach {
			//if ((Distance (SenderPosition, CurrentPlayer->LastPos) > m_PlayerIsOutOfReach)
			//&&  (CurrentPlayer->Callsign.compare (0, 3, "obs", 3) != 0))
			//{
			//CurrentPlayer++ 
			continue
		} 
		
		//  only send packet to local clients
		if loopPlayer.IsLocal {
			//SendChatMessages (CurrentPlayer);
			//m_DataSocket->sendto (Msg, Bytes, 0, &CurrentPlayer->Address);
			_, err := me.DataSocket.WriteToUDP(xdr_bytes, loopPlayer.Address)
			if err != nil {
				// TODO ?
			}
			loopPlayer.PktsSentTo++
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
	SendingPlayer := NewFG_Player() // placleholder 
	me.SendToRelays (xdr_bytes, length, SendingPlayer)
	
} // FgServer::HandlePacket ( char* sMsg[MAX_PACKET_SIZE] )



func (me *FgServer) DEADPacketIsValid(length int, header message.HeaderMsg, SenderAddress *net.UDPAddr ) bool {

	var ErrorMsg string

	// Check header Packet size
	s := int(unsafe.Sizeof(header))
	if length <  s {
		ErrorMsg  = SenderAddress.String()
		ErrorMsg += " packet size is too small!"
		fmt.Println("ERROR: PacketIsValid()", ErrorMsg)
		me.AddBadClient(SenderAddress, ErrorMsg, true)
		return false
	}
	
	//= Check magic
	if header.Magic != message.MSG_MAGIC && header.Magic != message.RELAY_MAGIC {
		ErrorMsg  = SenderAddress.String();
		ErrorMsg += " BAD magic number: "
		//ErrorMsg += MsgHdr.Magic // TODO
		//fmt.Println("TODO: Handle Wrong Magic")
		fmt.Println("ERROR: PacketIsValid()", ErrorMsg)
		me.AddBadClient(SenderAddress, ErrorMsg, true)
		return false
	}
	
	// Check Protocol Version
	//if MsgHdr.Version != message.PROTOCOL_VER {
	//	ErrorMsg  = SenderAddress.String()
	//	ErrorMsg += " BAD protocol version! Should be "
		// TODO bitshift
		//converter*    tmp;
		//tmp = (converter*) (& PROTO_VER);
		//ErrorMsg += NumToStr (tmp->High, 0);
		//ErrorMsg += "." + NumToStr (tmp->Low, 0);
		//ErrorMsg += " but is ";
		//tmp = (converter*) (& MsgHdr->Version);
		//ErrorMsg += NumToStr (tmp->Low, 0);
		//ErrorMsg += "." + NumToStr (tmp->High, 0);
	//	fmt.Println("ERROR: PacketIsValid()", ErrorMsg)
	//	me.AddBadClient(SenderAddress, ErrorMsg, true);
	//	return false
	//}
	/*
	if MsgHdr.Type == message.TYPE_POS {
		lenny := uint32( unsafe.Sizeof(&message.HeaderMsg) + unsafe.Sizeof(&message.PositionMsg{}) )
		if MsgHdr < lenny {
			ErrorMsg  = SenderAddress.String()
			ErrorMsg += " Client sends insufficient position data, "
			ErrorMsg += fmt.Sprintf( "should be %d", lenny)
			ErrorMsg += fmt.Sprintf(" is: %d", MsgHdr.MsgLen)
			me.AddBadClient (SenderAddress, ErrorMsg, true);
			return false
		}
	}*/
	return true
} // FgServer::PacketIsValid ()
