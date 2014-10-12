package flightgear

import (
	//"fmt"
)

import (
	"github.com/FreeFlightSim/go-fgms/simgear"
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



// Message Types
const (
	CHAT_MSG_ID = 1 //= is this used ??
	RESET_DATA_ID = 6
	POS_DATA_ID = 7
)



/* 
	XDR demands 4 byte alignment, but some compilers use 8 byte alignment
	so it's safe to let the overall size of a network message be a 
	multiple of 8!
*/
const ( 
	MAX_CALLSIGN_LEN	= 8
	MAX_CHAT_MSG_LEN   	= 256
	MAX_MODEL_NAME_LEN 	= 96
	MAX_PROPERTY_LEN   	= 52
)


/* T_MsgHdr - The Header At start of all messages sent 
	
	Original: http://gitorious.org/fgms/fgms-0-x/blobs/master/src/flightgear/MultiPlayer/mpmessages.hxx#line62
*/
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
	 CallsignBytes [MAX_CALLSIGN_LEN]byte //Callsign[MAX_CALLSIGN_LEN]
}

// Returns the Callsign as a String
func (me *T_MsgHdr) Callsign() string{
	return string(me.CallsignBytes[:])

}



// T_ChatMsg - A Chat message
type T_ChatMsg struct {
    Text [MAX_CHAT_MSG_LEN]byte
}


/* T_PositionMsg - Position Message
	
	Original Source: http://gitorious.org/fgms/fgms-0-x/blobs/master/src/flightgear/MultiPlayer/mpmessages.hxx#line78
	Note:
		all the important values are float32 
		with the exception of position which is float64
		This caused a clash with Point3D which needs to be either 32 or 64
		- For now the 32's are converted to 64's
*/
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






// Represents a Property Message 
type T_PropertyMsg struct{
    //xdr_data_t id;
    //xdr_data_t value;
    id uint32
    value uint32
}


  
type FGFloatPropertyData struct{
  //unsigned id;
  id uint32
  //float value;
  value float32
}


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

