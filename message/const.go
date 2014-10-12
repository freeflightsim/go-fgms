
package message

import (
)

// Magic value for messages - currently FGFS
const MSG_MAGIC = 0x46474653  // "FGFS"
//const MSG_MAGIC = "SFGF"  // "FGFS"



// Protocol Version - currently 1.1
const PROTOCOL_VER = 0x00010001  // 1.1

// Message Types
const (
	TYPE_CHAT = 1 //= is this used ??
	TYPE_RESET = 6
	TYPE_POS = 7
)



/*
	XDR demands Id4 byte alignment, but some compilers use 8 byte alignment
	so it's safe to let the overall size of a network message be a
	multiple of 8!
*/
const (
	MAX_CALLSIGN_LEN	= 8
	MAX_CHAT_MSG_LEN   	= 256
	MAX_MODEL_NAME_LEN 	= 96
	MAX_PROPERTY_LEN   	= 52
)
