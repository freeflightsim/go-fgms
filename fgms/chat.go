
package fgms


import (
	"fmt"
)

/**
 * @brief Create a chat message and put it into the internal message queue
 * @param ID int with the ?
 * @param Msg String with the message
 */
func (me *FG_SERVER) CreateChatMessage(ID int,Msg string){

	fmt.Println("CreateChatMessage", ID, Msg)
  /* T_MsgHdr        MsgHdr;
  T_ChatMsg       ChatMsg;
  unsigned int    NextBlockPosition = 0;
  char*           Message;
  int             len =  sizeof(T_MsgHdr) + sizeof(T_ChatMsg);

  MsgHdr.Magic            = XDR_encode<uint32_t> (MSG_MAGIC);
  MsgHdr.Version          = XDR_encode<uint32_t> (PROTO_VER);
  MsgHdr.MsgId            = XDR_encode<uint32_t> (CHAT_MSG_ID);
  MsgHdr.MsgLen           = XDR_encode<uint32_t> (len);
  MsgHdr.ReplyAddress     = 0;
  MsgHdr.ReplyPort        = XDR_encode<uint32_t> (m_ListenPort);
  strncpy(MsgHdr.Callsign, "*FGMS*", MAX_CALLSIGN_LEN);
  MsgHdr.Callsign[MAX_CALLSIGN_LEN - 1] = '\0';
  */
  /*
  while (NextBlockPosition < Msg.length())
  {
    strncpy (ChatMsg.Text, 
    Msg.substr (NextBlockPosition, MAX_CHAT_MSG_LEN - 1).c_str(),
    MAX_CHAT_MSG_LEN);
    ChatMsg.Text[MAX_CHAT_MSG_LEN - 1] = '\0';
    Message = new char[len];
    memcpy (Message, &MsgHdr, sizeof(T_MsgHdr));
    memcpy (Message + sizeof(T_MsgHdr), &ChatMsg,
    sizeof(T_ChatMsg));
    m_MessageList.push_back (mT_ChatMsg(ID,Message));
    NextBlockPosition += MAX_CHAT_MSG_LEN - 1;
  }*/
  
} // FG_SERVER::CreateChatMessage ()

