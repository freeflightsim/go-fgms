
package fgms

// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_server.cxx#line167
type FG_Server struct {

  /*typedef union
  {
    uint32_t    complete;
    int16_t     High;
    int16_t     Low;
  } converter; 
  converter*    tmp; */
  Initialized bool
  ReinitData bool
  ReinitTelnet bool
  ListenPort int
  PlayerExpires int
  Listening false
  //Loglevel            = SG_INFO;
  DataSocket int
  TelnetPort int
  NumMaxClients int
  PlayerIsOutOfReach int // nautical miles
  NumCurrentClients int
  IsParent bool
  MaxClientID int
  ServerName  string
  BindAddress string
  //tmp                   = (converter*) (& PROTO_VER);
  //ProtoMinorVersion   = tmp->High;
  //ProtoMajorVersion   = tmp->Low;
  //LogFileName         = DEF_SERVER_LOG; // "fg_server.log";
  //wp                  = fopen("wp.txt", "w");
  BlackList           = map<uint32_t, bool>();
  RelayMap            = map<uint32_t, string>();
  IsTracked bool
  Tracker int
  //UpdateSecs          = DEF_UPDATE_SECS;
  // clear stats - should show what type of packet was received
  PacketsReceived     = 0;
  TelnetReceived      = 0;
  BlackRejected       = 0;  // in black list
  PacketsInvalid      = 0;  // invalid packet
  UnknownRelay        = 0;  // unknown relay
  RelayMagic          = 0;  // relay magic packet
  PositionData        = 0;  // position data packet
  NotPosData          = 0;
  // clear totals
  MT_PacketsReceived int
  MT_BlackRejected int
  MT_PacketsInvalid int
  MT_UnknownRelay int
  MT_PositionData int
  MT_TelnetReceived int
  MT_RelayMagic int
  MT_NotPosData int
  
  CrossFeedFailed int
  CrossFeedSent int
  
  MT_CrossFeedFailed int
  MT_CrossFeedSent int
  TrackerConnect int
  TrackerDisconnect int
  TrackerPostion int // Tracker messages queued
  //pthread_mutex_init( &m_PlayerMutex, 0 );
} 



