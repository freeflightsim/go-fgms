
package fgms


import (
	"fmt"
	"unsafe"
)
import(
	"github.com/fgx/go-fgms/flightgear"
)

/*  Create a chat message and put it into the internal message queue

	This needs completing
*/
func (me *FG_SERVER) CreateChatMessage(ID int, Msg string){

	fmt.Println("CreateChatMessage", ID, Msg)
	//T_MsgHdr        MsgHdr;
	//T_ChatMsg       ChatMsg;
	//unsigned int    NextBlockPosition = 0;
	//char*           Message;
	//int             len =  sizeof(T_MsgHdr) + sizeof(T_ChatMsg);
	
	var MsgHdr flightgear.T_MsgHdr
	var ChatMsg flightgear.T_ChatMsg
	
	var NextBlockPosition uint  = 0
	var Message []byte //char*           Message;
	
	var lenny int =  int( unsafe.Sizeof(MsgHdr) + unsafe.Sizeof(ChatMsg) )
	
	fmt.Println(NextBlockPosition, lenny, Message)
	
	MsgHdr.Magic            = flightgear.MSG_MAGIC
	MsgHdr.Version          = flightgear.PROTO_VER
	MsgHdr.MsgId            = flightgear.CHAT_MSG_ID
	//MsgHdr.MsgLen           = XDR_encode<uint32_t> (len);
	
	// Depreciated
	MsgHdr.ReplyAddress     = 0
	MsgHdr.ReplyPort        = uint32(me.ListenPort)
	
	//MsgHdr.Callsign = [8]byte ("*FGMS*") + [0] + [0]
	//strncpy(MsgHdr.Callsign, "*FGMS*", MAX_CALLSIGN_LEN);
	//MsgHdr.Callsign[MAX_CALLSIGN_LEN - 1] = '\0';
	//	
	
	// MsgHdr.Callsign is  Callsign [8]byte 
	// There's got to be an easier way to do this in GO!
	cs_bytes := [8]byte{0,0,0,0,0,0,0,0} 
	for idx, char := range("*FGMS*") {
	 	cs_bytes[idx] = byte(char)
	}
	MsgHdr.Callsign = cs_bytes 
	

	
	//while (NextBlockPosition < Msg.length())
	//{
		//strncpy (ChatMsg.Text, 
		//Msg.substr (NextBlockPosition, MAX_CHAT_MSG_LEN - 1).c_str(),
		//MAX_CHAT_MSG_LEN);
		//ChatMsg.Text[MAX_CHAT_MSG_LEN - 1] = '\0';
		//Message = new char[len];
		//memcpy (Message, &MsgHdr, sizeof(T_MsgHdr));
		//memcpy (Message + sizeof(T_MsgHdr), &ChatMsg,
		//sizeof(T_ChatMsg));
		//m_MessageList.push_back (mT_ChatMsg(ID,Message));
		//NextBlockPosition += MAX_CHAT_MSG_LEN - 1;
	//}
	//while (NextBlockPosition < Msg.length())
	//{
	var idx uint = 0
	for x, cha := range Msg {
		
		ChatMsg.Text[x] = byte(cha)
		
		idx++
		
		if idx == flightgear.MAX_CHAT_MSG_LEN - 1 {
			// this message is too long so send this part ?
			me.MessageList = append(me.MessageList, ChatMsg)
		}
		//strncpy (ChatMsg.Text, 
		           //Msg.substr (NextBlockPosition, MAX_CHAT_MSG_LEN - 1).c_str(),
		           //MAX_CHAT_MSG_LEN);
		//ChatMsg.Text[MAX_CHAT_MSG_LEN - 1] = '\0';
		//Message = new char[len];
		//memcpy (Message, &MsgHdr, sizeof(T_MsgHdr));
		//memcpy (Message + sizeof(T_MsgHdr), &ChatMsg,
		//sizeof(T_ChatMsg));
		//m_MessageList.push_back (mT_ChatMsg(ID,Message));
		//NextBlockPosition += MAX_CHAT_MSG_LEN - 1;
	}
	// We got a message to send anyway ?
	me.MessageList = append(me.MessageList, ChatMsg)
	
} // FG_SERVER::CreateChatMessage ()




/* Send any message in m_MessageList to client
	 @param CurrentPlayer Player to send message to
*/
func (me *FG_SERVER) SendChatMessages() {

	//mT_MessageIt  CurrentMessage;
		/*
	if ((CurrentPlayer->IsLocal) && (m_MessageList.size()))
	{
		CurrentMessage = m_MessageList.begin();
		while (CurrentMessage != m_MessageList.end())
		{
		if ((CurrentMessage->Target == 0)
		||  (CurrentMessage->Target == CurrentPlayer->ClientID))
		{
			int len = sizeof(T_MsgHdr) + sizeof(T_ChatMsg);
			m_DataSocket->sendto (CurrentMessage->Msg, len, 0,
			&CurrentPlayer->Address);
		}
		CurrentMessage++;
		}
	} */
	
} // FG_SERVER::SendChatMessages ()
