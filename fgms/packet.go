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
	"github.com/davecgh/go-xdr/xdr"

	"github.com/fgx/go-fgms/flightgear"
)


//------------------------------------------------------------------------

// Handle client connections
func (me *FG_SERVER) HandlePacket(Msg []byte, Bytes int, SenderAddress *net.UDPAddr){
	
	//T_MsgHdr*       MsgHdr;
	var MsgHdr flightgear.T_MsgHdr
	//T_PositionMsg*  PosMsg;
	var PosMsg flightgear.T_PositionMsg
	
	//uint32_t        MsgId;
	//uint32_t        MsgMagic;
	//Timestamp time.Time
	Timestamp := Now()
	
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
	
	remainingBytes, err := xdr.Unmarshal(Msg, &MsgHdr)
	if err != nil{
		fmt.Println("XDR Decode Error", err)
		return
	}
	fmt.Println("remain=", len(remainingBytes), SenderAddress)
	
	//MsgMagic  = XDR_decode<uint32_t> (MsgHdr->Magic);
	//MsgId     = XDR_decode<uint32_t> (MsgHdr->MsgId);
	//fmt.Println( "Magic/ID", MsgHdr.Magic, MsgHdr.Version, MsgHdr.MsgId, MsgHdr.Callsign, MsgHdr.ReplyAddress, MsgHdr.ReplyPort )
	
	//fmt.Println ("=magic=", flightgear.MSG_MAGIC == MsgHdr.Magic) //WORKS
	//fmt.Println ("=proto=", flightgear.PROTO_VER == MsgHdr.Version) //WORKS
	//fmt.Println ("=ID=", MsgHdr.MsgId)
	//cs := "" //string(MsgHdr.Callsign[0]) + string(MsgHdr.Callsign[1]) + string(MsgHdr.Callsign[2]) + string(MsgHdr.Callsign[3]) + string(MsgHdr.Callsign[4]) + string(MsgHdr.Callsign[5]) + string(MsgHdr.Callsign[6]) + string(MsgHdr.Callsign[7])
	//for _, ele := range MsgHdr.Callsign{
	//	if ele != 0 {
	//		cs += string(ele)
	//	}
	//}    
	fmt.Println ("=Got Header=", MsgHdr.Callsign, MsgHdr.CallsignString())
	

	//------------------------------------------------------
	// First of all, send packet to all crossfeed servers.
	//SendToCrossfeed (Msg, Bytes, SenderAddress); ?? SHould then be send pre vaildation ?
	me.SendToCrossfeed(Msg, Bytes, SenderAddress)


	//------------------------------------------------------
	//=  Now do the local processing TODO
	//if me.IsBlackListed(SenderAddress) {
	//	me.BlackListRejected++
	//	return
	//}
	
	// Check packet is valid
	fmt.Println (" > checkvalid")
	if !me.PacketIsValid(Bytes, MsgHdr, SenderAddress) {
		me.PacketsInvalid++
		fmt.Println ("  <<  NO checkvalid")
		return
	} 
	fmt.Println ("  <<  YES checkvalid")
	
	if MsgHdr.Magic == RELAY_MAGIC { // not a local client
		if !me.IsKnownRelay(SenderAddress) {
			me.UnknownRelay++ 
			return
		}else{
			me.RelayMagic++ // bump relay magic packet
		}
	}
	
	//////////////////////////////////////////////////
	//    Store senders position
	//////////////////////////////////////////////////
	if MsgHdr.MsgId == flightgear.POS_DATA_ID	{
		me.PositionData++
		remainingBytes2, errPos := xdr.Unmarshal(remainingBytes, &PosMsg)
		if err != nil{
			fmt.Println("XDR Decode Position Error", errPos)
			return
		}
		fmt.Println("remain2=", len(remainingBytes2), PosMsg.Model)
	
		//PosMsg = (T_PositionMsg *) (Msg + sizeof(T_MsgHdr));
		//double x = XDR_decode64<double> (PosMsg->position[X]);
		//double y = XDR_decode64<double> (PosMsg->position[Y]);
		//double z = XDR_decode64<double> (PosMsg->position[Z]);
		x := PosMsg.Position[X]
		y := PosMsg.Position[Y]
		z := PosMsg.Position[Z]
		if x == 0.0 || y == 0.0 || z == 0.0 { // ignore while position is not settled
			return
		}
		SenderPosition.Set (x, y, z);
		
		/* SenderOrientation.Set (
			XDR_decode<float> (PosMsg->orientation[X]),
			XDR_decode<float> (PosMsg->orientation[Y]),
			XDR_decode<float> (PosMsg->orientation[Z])
		)*/
		//TODO Wrong TYPE wtf!
		//SenderOrientation.Set(PosMsg.Orientation[X], PosMsg.Orientation[Y],	PosMsg.Orientation[Z])
		SenderOrientation.Set(0,0,0)
	} else {
		me.NotPosData++
	} 
	
	// Add Client to list if its not known
	senderInList := me.SenderIsKnown(MsgHdr.CallsignString())
	fmt.Println ("  <<  senderInList", senderInList)
	if senderInList == SENDER_UNKNOWN { 
		// unknown, add to the list
		if MsgHdr.MsgId != flightgear.POS_DATA_ID {
			return // ignore client until we have a valid position
		}
		//tempPosMsg := flightgear.T_PositionMsg{}
		me.AddClient(SenderAddress, MsgHdr, PosMsg)
		
	}else if senderInList == SENDER_DIFF_IP {
		return // known, but different IP => ignore
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
	xCallsign := MsgHdr.CallsignString()
	xIsObserver :=  strings.ToLower(MsgHdr.CallsignString())[0:3] ==  "obs"
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
			if MsgHdr.MsgId == flightgear.POS_DATA_ID 	{
				loopPlayer.LastPos         = SenderPosition
				loopPlayer.LastOrientation = SenderOrientation
			}else{
				SenderPosition    = loopPlayer.LastPos
				SenderOrientation = loopPlayer.LastOrientation
			}
			//SendingPlayer = CurrentPlayer
			loopPlayer.Timestamp = Timestamp
			loopPlayer.PktsReceivedFrom++
			//CurrentPlayer++;
			continue; // don't send packet back to sender
		}
		///     do not send packets to clients if the
		//      origin is an observer, but do send
		//      chat messages anyway
		//      FIXME: MAGIC = SFGF!
		if xIsObserver && MsgHdr.MsgId != flightgear.CHAT_MSG_ID {
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
			_, err := me.DataSocket.WriteToUDP(Msg, loopPlayer.Address)
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
	me.SendToRelays (Msg, Bytes, SendingPlayer)
	
} // FG_SERVER::HandlePacket ( char* sMsg[MAX_PACKET_SIZE] )



func (me *FG_SERVER) PacketIsValid(	Bytes int, MsgHdr flightgear.T_MsgHdr, SenderAddress *net.UDPAddr ) bool {

	var ErrorMsg string

	// Check header Packet size
	s := int(unsafe.Sizeof(MsgHdr))
	if Bytes <  s {
		ErrorMsg  = SenderAddress.String()
		ErrorMsg += " packet size is too small!"
		fmt.Println("ERROR: PacketIsValid()", ErrorMsg)
		me.AddBadClient(SenderAddress, ErrorMsg, true)
		return false
	}
	
	//= Check magic
	if MsgHdr.Magic != flightgear.MSG_MAGIC && MsgHdr.Magic != RELAY_MAGIC {
		ErrorMsg  = SenderAddress.String();
		ErrorMsg += " BAD magic number: "
		//ErrorMsg += MsgHdr.Magic // TODO
		//fmt.Println("TODO: Handle Wrong Magic")
		fmt.Println("ERROR: PacketIsValid()", ErrorMsg)
		me.AddBadClient(SenderAddress, ErrorMsg, true)
		return false
	}
	
	// Check Protocol Version
	if MsgHdr.Version != flightgear.PROTO_VER {
		ErrorMsg  = SenderAddress.String()
		ErrorMsg += " BAD protocol version! Should be "
		// TODO bitshift
		//converter*    tmp;
		//tmp = (converter*) (& PROTO_VER);
		//ErrorMsg += NumToStr (tmp->High, 0);
		//ErrorMsg += "." + NumToStr (tmp->Low, 0);
		//ErrorMsg += " but is ";
		//tmp = (converter*) (& MsgHdr->Version);
		//ErrorMsg += NumToStr (tmp->Low, 0);
		//ErrorMsg += "." + NumToStr (tmp->High, 0);
		fmt.Println("ERROR: PacketIsValid()", ErrorMsg)
		me.AddBadClient(SenderAddress, ErrorMsg, true);
		return false
	} 
	
	if MsgHdr.MsgId == flightgear.POS_DATA_ID {
		lenny := uint32( unsafe.Sizeof(MsgHdr) + unsafe.Sizeof(&flightgear.T_PositionMsg{}) )
		if MsgHdr.MsgLen < lenny {
			ErrorMsg  = SenderAddress.String()
			ErrorMsg += " Client sends insufficient position data, "
			ErrorMsg += fmt.Sprintf( "should be %d", lenny)
			ErrorMsg += fmt.Sprintf(" is: %d", MsgHdr.MsgLen)
			me.AddBadClient (SenderAddress, ErrorMsg, true);
			return false
		}
	}
	return true
} // FG_SERVER::PacketIsValid ()
