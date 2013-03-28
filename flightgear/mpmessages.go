package flightgear

import (
	"fmt"
	//"path/filepath"
)

import (
	"github.com/fgx/go-fgms/simgear"
)

/*
	This section is used as the Protocol. .for more info
	- Original Sources between
		fgms: http://gitorious.org/fgms/fgms-0-x/blobs/master/src/flightgear/MultiPlayer/mpmessages.hxx
		fg:  http://gitorious.org/fg/flightgear/trees/next/src/MultiPlayer/
	- INFO: http://wiki.flightgear.org/Multiplayer_protocol
		
	Note the xdr part is currently in the fgms and not here. Although that is a possible
	- This lib external = http://godoc.org/github.com/davecgh/go-xdr/xdr
*/	

// Magic value for messages - currently FGFS
const MSG_MAGIC = 0x46474653  // "FGFS"


// Protocol Version - currently 1.1
const PROTO_VER = 0x00010001  // 1.1

// pete FAIL FAIL FAIL's 
func GetProtocolVerString() string {
	return "1.1"
	major := PROTO_VER >> 16
	minor := PROTO_VER >> 4
	return fmt.Sprintf("%d.%d", major, minor) 
}


// Message Types
const (
	CHAT_MSG_ID = 1
	RESET_DATA_ID = 6
	POS_DATA_ID = 7
)



/* 
	XDR demands 4 byte alignment, but some compilers use8 byte alignment
	so it's safe to let the overall size of a network message be a 
	multiple of 8!
*/
const ( 
	MAX_CALLSIGN_LEN 	 = 8
	MAX_CHAT_MSG_LEN   = 256
	MAX_MODEL_NAME_LEN = 96
	MAX_PROPERTY_LEN   = 52
)


// T_MsgHdr - Header for use with all messages sent 
// typedef uint32_t    xdr_data_t;      /* 4 Bytes */
// typedef uint64_t    xdr_data2_t;     /* 8 Bytes */
 type T_MsgHdr struct {
	
	// Magic Value
    Magic uint32 //xdr_data_t
    
    // Protocol version
    Version uint32 //xdr_data_t            
    
    // Message identifier 
    MsgId uint32 //xdr_data_t    
    
    // Absolute length of message
    MsgLen uint32 //xdr_data_t   
    
    // DEPRECEATED: Player's receiver address 
    ReplyAddress uint32 //xdr_data_t   
    
    // DEPRECEATED: Player's receiver port
    ReplyPort uint32 //xdr_data_t   
    
    // Callsign used by the player 
    Callsign [MAX_CALLSIGN_LEN]byte //Callsign[MAX_CALLSIGN_LEN] 
}

// TODO There has Got to be a better way
func (me *T_MsgHdr) CallsignString() string{
	s := ""
	for _, ele := range me.Callsign {
		if ele == 0 {
			return s
		}
		s += string(ele)
	}
	return s   
}



// T_ChatMsg - Chat message
type T_ChatMsg struct {
	
	// Text of chat message 
    //string Text //char Text[MAX_CHAT_MSG_LEN];  
    Text [MAX_CHAT_MSG_LEN]byte
}


// T_PositionMsg - Position Message
type T_PositionMsg struct{
	
	/// Name of the aircraft model 
    // - char Model[MAX_MODEL_NAME_LEN]; 
    Model [MAX_MODEL_NAME_LEN]byte

    // Time when this packet was generated
    // - xdr_data2_t time;
    Time uint64
	
	/// Time offset for network lag ? 
    // - xdr_data2_t lag;
    Lag uint64

    // Position wrt the earth centered frame
    // - xdr_data2_t position[3];
    Position [3]float64
	
	
    // Orientation wrt the earth centered frame, stored in the angle axis
    // representation where the angle is coded into the axis length
    // - xdr_data_t orientation[3];
    Orientation [3]float32 //uint32

	// Linear velocity wrt the earth centered frame measured in the earth centered frame
    // - xdr_data_t linearVel[3];
    LinearVel [3]float32 //uint32
	
    // Angular velocity wrt the earth centered frame measured in the earth centered frame
    // - xdr_data_t angularVel[3];
    AngularVel [3]float32 // uint32

	// Linear acceleration wrt the earth centered frame measured in the earth centered frame
    // - xdr_data_t linearAccel[3];
    LinearAccel [3]float32 // uint32
	
    // Angular acceleration wrt the earth centered frame measured in the earth centered frame
    // - xdr_data_t angularAccel[3];
    AngularAccel [3]float32 //uint32
}

// Returns the Model as a string 
// - TODO There has Got to be a better way
func (me *T_PositionMsg) ModelString() string{
	s := ""
	for _, ele := range me.Model {
		if ele == 0 {
			return s
		}
		s += string(ele)
	}
	return s   
}






// Represents a Property Message for XDR 
type T_PropertyMsg struct{
    //xdr_data_t id;
    //xdr_data_t value;
    id uint32
    value uint32
}


/**
 * @struct FGFloatPropertyData  
 * @brief Property Data 
 */
type FGFloatPropertyData struct{
  //unsigned id;
  id uint32
  //float value;
  value float32
}

/** @brief Position Message */
type FGExternalMotionData struct {
	
  /** 
   * @brief Simulation time when this packet was generated 
   */
  //double time;
  time uint64
  
  	/*
   	The artificial lag the client should stay behind the average
   	simulation time to arrival time diference
   	
	todo -  should be some 'per model' instead of 'per packet' property  double lag;
   	Position wrt the earth centered frame  
  	
  	- SGVec3d position
  	*/
  	position simgear.SGVec3d 
  
  	// Orientation wrt the earth centered frame 
  	//SGQuatf orientation;
	orientation simgear.SGQuatf
  
  	// Linear velocity wrt the earth centered frame measured in
	//the earth centered frame
	// - SGVec3f linearVel;
  	linearVel simgear.SGVec3f
  
  	// Angular velocity wrt the earth centered frame measured in the earth centered frame
  	// - SGVec3f angularVel;
  	angularVel simgear.SGVec3f
  
  	// Linear acceleration wrt the earth centered frame measured in the earth centered frame
  	// - SGVec3f linearAccel;
  	linearAccel simgear.SGVec3f
   
  	// Angular acceleration wrt the earth centered frame measured in the earth centered frame 
  	// - SGVec3f angularAccel;
  	angularAccel simgear.SGVec3f
  
  	// The set of properties recieved for this timeslot 
  	// TODO std::vector<FGFloatPropertyData> properties;
}

