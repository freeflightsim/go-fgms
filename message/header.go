
package message

type HeaderMsg struct {

	// Magic Value
	Magic uint32 //xdr_data_t

	// Protocol version
	Version uint32 //xdr_data_t

	// Message identifier
	Type uint32 //xdr_data_t

	// Absolute length of message
	Len uint32 //xdr_data_t

	// DEPRECEATED: Player's receiver address
	ReplyAddress uint32 //xdr_data_t

	// DEPRECEATED: Player's receiver port
	ReplyPort uint32 //xdr_data_t

	// Callsign used by the player
	CallsignBytes [MAX_CALLSIGN_LEN]byte //Callsign[MAX_CALLSIGN_LEN]
}

// returns Callsign as string
func (me *HeaderMsg) Callsign() string{
	return string(me.CallsignBytes[:])
}
