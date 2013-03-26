package flightgear

import (
	"fmt"
)

import (
	"github.com/fgx/go-fgms/simgear"
	//"github.com/davecgh/go-xdr/xdr"
)

// TODO:  XDR decoding should be here ?

// magic value for messages 
const MSG_MAGIC = 0x46474653  // "FGFS"


//   protocol version 
const PROTO_VER = 0x00010002  // 1.1
const PROTO_VER_STR = "1.1"

func GetProtocolVerString() string {
	major := PROTO_VER >> 16
	minor := PROTO_VER >> 4
	return fmt.Sprintf("%d.%d", major, minor) // FAIL FAIL FAIL 
}


const CHAT_MSG_ID = 1
const RESET_DATA_ID = 6
const POS_DATA_ID = 7


/** @brief Internal Constants */
//enum FG_SERVER_CONSTANTS
//{
// return values

	

/* 
XDR demands 4 byte alignment, but some compilers use8 byte alignment
so it's safe to let the overall size of a network message be a 
multiple of 8!
*/
const MAX_CALLSIGN_LEN 	 = 8
const MAX_CHAT_MSG_LEN   = 256
const MAX_MODEL_NAME_LEN = 96
const MAX_PROPERTY_LEN   = 52

// External = http://godoc.org/github.com/davecgh/go-xdr/xdr


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
    
    //  Player's receiver address 
    ReplyAddress uint32 //xdr_data_t   
    
    // Player's receiver port
    ReplyPort uint32 //xdr_data_t   
    
    /// Callsign used by the player 
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
    Text [MAX_CALLSIGN_LEN]byte
}


// T_PositionMsg - Position Message
type T_PositionMsg struct{
	
	/** @brief  Name of the aircraft model */
    //char Model[MAX_MODEL_NAME_LEN]; 
    Model [MAX_MODEL_NAME_LEN]byte

    /** @brief Time when this packet was generated */
    //xdr_data2_t time;
    time uint64
	
	/** @brief Time when this packet was generated */
    //xdr_data2_t lag;
    lag uint64

    /** @brief Position wrt the earth centered frame */
    //xdr_data2_t position[3];
    position [3]uint64
	
	
    /** @brief Orientation wrt the earth centered frame, stored in the angle axis
     *         representation where the angle is coded into the axis length
	 */
    //xdr_data_t orientation[3];
    orientation [3]uint32

	/** @brief Linear velocity wrt the earth centered frame measured in
     *         the earth centered frame
	 */
    //xdr_data_t linearVel[3];
    linearVel [3]uint32
	
    /** @brief Angular velocity wrt the earth centered frame measured in
     *          the earth centered frame
	 */
    //xdr_data_t angularVel[3];
    angularVel [3]uint32

	/** @brief Linear acceleration wrt the earth centered frame measured in
     *         the earth centered frame
	 */
    //xdr_data_t linearAccel[3];
    linearAccel [3]uint32
	
    /** @brief Angular acceleration wrt the earth centered frame measured in
     *         the earth centered frame
	 */
    //xdr_data_t angularAccel[3];
    angularAccel [3]uint32
}

// TODO There has Got to be a better way
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


/** 
 * @struct T_PropertyMsg 
 *  @brief Property Message 
 */
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

