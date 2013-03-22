package mpserver


import (

)

// This is ported from an for research 
// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/flightgear/MultiPlayer/mpmessages.hxx


// magic value for messages 
const MSG_MAGIC = 0x46474653  // "FGFS"


//   protocol version 
const PROTO_VER = 0x00010001;  // 1.1


const CHAT_MSG_ID = 1
const RESET_DATA_ID = 6
const POS_DATA_ID = 7

/* 
XDR demands 4 byte alignment, but some compilers use8 byte alignment
so it's safe to let the overall size of a network message be a 
multiple of 8!
*/
const MAX_CALLSIGN_LEN 	 = 8
const MAX_CHAT_MSG_LEN   = 256
const MAX_MODEL_NAME_LEN = 96
const MAX_PROPERTY_LEN   = 52


// T_MsgHdr - Header for use with all messages sent 
 struct T_MsgHdr {
	
	// Magic Value
    Magic xdr_data_t
    
    // Protocol version
    Version xdr_data_t            
    
    // Message identifier 
    MsgId xdr_data_t    
    
    // Absolute length of message
    MsgLen xdr_data_t   
    
    //  Player's receiver address 
    ReplyAddress xdr_data_t   
    
    // Player's receiver port
    ReplyPort xdr_data_t   
    
    /// Callsign used by the player 
    string Callsign //Callsign[MAX_CALLSIGN_LEN] 
}


// T_ChatMsg - Chat message
struct T_ChatMsg {
	
	// Text of chat message 
    string Text //char Text[MAX_CHAT_MSG_LEN];  
}


/** 
 * @struct T_PositionMsg
 * @brief Position Message
 */
struct T_PositionMsg {
	
	/** @brief  Name of the aircraft model */
    char Model[MAX_MODEL_NAME_LEN]; 

    /** @brief Time when this packet was generated */
    xdr_data2_t time;
	
	/** @brief Time when this packet was generated */
    xdr_data2_t lag;

    /** @brief Position wrt the earth centered frame */
    xdr_data2_t position[3];
	
	
    /** @brief Orientation wrt the earth centered frame, stored in the angle axis
     *         representation where the angle is coded into the axis length
	 */
    xdr_data_t orientation[3];

	/** @brief Linear velocity wrt the earth centered frame measured in
     *         the earth centered frame
	 */
    xdr_data_t linearVel[3];
	
    /** @brief Angular velocity wrt the earth centered frame measured in
     *          the earth centered frame
	 */
    xdr_data_t angularVel[3];

	/** @brief Linear acceleration wrt the earth centered frame measured in
     *         the earth centered frame
	 */
    xdr_data_t linearAccel[3];
	
    /** @brief Angular acceleration wrt the earth centered frame measured in
     *         the earth centered frame
	 */
    xdr_data_t angularAccel[3];
};


/** 
 * @struct T_PropertyMsg 
 *  @brief Property Message 
 */
struct T_PropertyMsg {
    xdr_data_t id;
    xdr_data_t value;
};


/**
 * @struct FGFloatPropertyData  
 * @brief Property Data 
 */
struct FGFloatPropertyData {
  unsigned id;
  float value;
};

/** @brief Position Message */
struct FGExternalMotionData {
	
  /** 
   * @brief Simulation time when this packet was generated 
   */
  double time;
  
  /** 
   * @brief The artificial lag the client should stay behind the average
   *        simulation time to arrival time diference
   * @todo  should be some 'per model' instead of 'per packet' property  double lag;
   *        Position wrt the earth centered frame  
   */
  SGVec3d position;
  
  /** @brief Orientation wrt the earth centered frame */
  SGQuatf orientation;
  
  /**
   * @brief Linear velocity wrt the earth centered frame measured in
   *        the earth centered frame
   */
  SGVec3f linearVel;
  
  /** 
   * @brief Angular velocity wrt the earth centered frame measured in the earth centered frame
   */
  SGVec3f angularVel;
  
  /** @brief Linear acceleration wrt the earth centered frame measured in the earth centered frame */
  SGVec3f linearAccel;
  
  /** @brief Angular acceleration wrt the earth centered frame measured in the earth centered frame */
  SGVec3f angularAccel;
  
  /** @brief The set of properties recieved for this timeslot */
  std::vector<FGFloatPropertyData> properties;
};

#endif
