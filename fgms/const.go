
package fgms



const (
	SUCCESS                 = 0
	ERROR_COMMANDLINE
	ERROR_CREATE_SOCKET
	ERROR_COULDNT_BIND
	ERROR_NOT_LISTENING
	ERROR_COULDNT_LISTEN
)

// other constants
const MAX_PACKET_SIZE         = 1024
const UPDATE_INACTIVE_PERIOD  = 1
const MAX_TELNETS             = 5



const (
	SENDER_UNKNOWN  = 0
	SENDER_KNOWN
	SENDER_DIFF_IP  // Not sure this is used
)
