
package fgms

import(
	"fmt"
)
import(
	"github.com/fgx/go-fgms/tracker"
)


const SUCCESS                 = 0
const ERROR_COMMANDLINE       = 1
const ERROR_CREATE_SOCKET     = 2
const ERROR_COULDNT_BIND      = 3
const ERROR_NOT_LISTENING     = 4
const ERROR_COULDNT_LISTEN    = 5

// other constants
const MAX_PACKET_SIZE         = 1024
const UPDATE_INACTIVE_PERIOD  = 1
const MAX_TELNETS             = 5
const RELAY_MAGIC             = 0x53464746    // GSGF




		
// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_server.cxx#line167
type FG_SERVER struct {

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
  Listening bool
  
  ServerName  string
  BindAddress string
  ListenPort int
  
  PlayerExpires int
  
  IamHUB bool
 
  //Loglevel            = SG_INFO;
  DataSocket int
  TelnetPort int
  NumMaxClients int
  PlayerIsOutOfReach int // nautical miles
  NumCurrentClients int
  IsParent bool
  MaxClientID int

	LogFileName string
	
  //tmp                   = (converter*) (& PROTO_VER);
  //ProtoMinorVersion   = tmp->High;
  //ProtoMajorVersion   = tmp->Low;
  //LogFileName         = DEF_SERVER_LOG; // "fg_server.log";
  //wp                  = fopen("wp.txt", "w");
  BlackList map[string]bool
  //RelayMap            = map<uint32_t, string>();
  
  IsTracked bool
  Tracker *tracker.FG_TRACKER
  
  //UpdateSecs          = DEF_UPDATE_SECS;
  // clear stats - should show what type of packet was received
  PacketsReceived int
  TelnetReceived int 
  BlackRejected int
  PacketsInvalid int //     = 0;  // invalid packet
  UnknownRelay int //       = 0;  // unknown relay
  RelayMagic  int //        = 0;  // relay magic packet
  PositionData int //         = 0;  // position data packet
  NotPosData int    //     = 0;
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

// Consruct and return pointer to new FG_SERVER instance
func NewFG_SERVER() *FG_SERVER {
	ob := new(FG_SERVER)
	ob.BlackList = make(map[string]bool)
	// set other defaults here
	return ob
}

func (me *FG_SERVER) SetServerName(name string){
	me.ServerName = name
}
func (me *FG_SERVER) SetBindAddress(addr string){
	me.BindAddress = addr
}

func (me *FG_SERVER) SetDataPort(port int){
	me.ListenPort = port
	me.ReinitData = true
}

func (me *FG_SERVER) SetTelnetPort(port int){
	me.TelnetPort = port
	me.ReinitTelnet = true
}



// Set nautical miles two players must be apart to be out of reach
func (me *FG_SERVER) SetOutOfReach(nm int){
	me.PlayerIsOutOfReach = nm
}

// Set time in seconds. if no packet arrives from a client
// within this time, the connection is dropped.  
func (me *FG_SERVER) SetPlayerExpires(secs int){
	me.PlayerExpires = secs
}

// Set if we are running as a Hubserver
func (me *FG_SERVER) SetHub(am_hub bool){
	me.IamHUB = am_hub
}


// Set the logfile
func (me *FG_SERVER) SetLogfile( log_file_name string){
	
	me.LogFileName = log_file_name
	
	/*TODO after research
  if (m_LogFile)
  {
    m_LogFile.close ();
  }
  m_LogFileName = LogfileName;
  m_LogFile.open (m_LogFileName.c_str(), ios::out|ios::app);
  sglog().enable_with_date (true);
  sglog().set_output (m_LogFile);
  */
} // FG_SERVER::SetLogfile ( const std::string &LogfileName )



// Insert a new relay server into internal list 
func (me *FG_SERVER) AddRelay(server string, port int) {
  //mT_Relay        NewRelay;
  //unsigned int    IP;
	NewRelay := NewMT_Relay(server, port)	
  	//NewRelay.Name = server
  	//NewRelay.Address.set ((char*) Server.c_str(), Port);
  	IP := NewRelay.Address.GetIP()
  	fmt.Println("New RELAY IP=", IP)
  	/*
  	if IP != INADDR_ANY && IP != INADDR_NONE {
    	m_RelayList.push_back (NewRelay);
    	string S; unsigned I;
   	 I = NewRelay.Name.find (".");
    if (I != string::npos)
    {
      S = NewRelay.Name.substr (0, I);
    }
    else
    {
      S = NewRelay.Name;
    }
    m_RelayMap[NewRelay.Address.getIP()] = S;
  	*/
} // FG_SERVER::AddRelay()




//////////////////////////////////////////////////////////////////////
/**
 * @brief Add a tracking server
 * @param Server String with server
 * @param Port The port number
 * @param IsTracked Is Stracked
 * @retval int -1 for fail or SUCCESS
 */
func (me *FG_SERVER) AddTracker(host string, port int, isTracked bool){
	me.IsTracked = isTracked
	me.Tracker = tracker.NewFG_TRACKER(host, port, 0)
	
	/* TODO
#ifndef NO_TRACKER_PORT
#ifdef USE_TRACKER_PORT
  if ( m_Tracker )
  {
    delete m_Tracker;
  }
  m_Tracker = new FG_TRACKER(Port,Server,0);
#else // !#ifdef USE_TRACKER_PORT
  if ( m_Tracker )
  {
    msgctl(m_ipcid,IPC_RMID,NULL);
    delete m_Tracker;
    m_Tracker = 0; // just deleted
  }
  printf("Establishing IPC\n");
  m_ipcid         = msgget(IPC_PRIVATE,IPCPERMS);
  if (m_ipcid <= 0)
  {
    perror("msgget getting ipc id failed");
    return -1;
  }
  m_Tracker = new FG_TRACKER(Port,Server,m_ipcid);
#endif // #ifdef USE_TRACKER_PORT y/n
#endif // NO_TRACKER_PORT
  return (SUCCESS);
  */
} // FG_SERVER::AddTracker()



// Add an IP to the blacklist
func (me *FG_SERVER) AddBlacklist(FourDottedIP string) {
  	//SG_ALERT (SG_SYSTEMS, SG_ALERT, "Adding to blacklist: " << FourDottedIP);
  	//m_BlackList[netAddress(FourDottedIP.c_str(), 0).getIP()] = true;
  	fmt.Println("Added to blacklist=", FourDottedIP)
  	me.BlackList[FourDottedIP] = true
} // FG_SERVER::AddBlacklist()
