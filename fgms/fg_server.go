
package fgms

import(
	"fmt"
	"log"
	"net"
)
import(
	"github.com/fgx/go-fgms/tracker"
	//"github.com/fgx/go-fgms/tracker"
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


type Version struct{
	Minor int
	Major int
}
		
// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_server.cxx#line167
type FG_SERVER struct {

  /*typedef union
  {
    uint32_t    complete;
    int16_t     High;
    int16_t     Low;
  } converter; 
  converter*    tmp; */
  VERSION int
  Initialized bool
  ReinitData bool
  ReinitTelnet bool
  Listening bool
  
  
	ServerName  string
  	BindAddress string
  	ListenPort int
  	IamHUB bool
  	
  	ServerVersion *Version
  	ProtocolVersion *Version
  
  	PlayerExpires int
  
	Telnet *TelnetServer
 
  //Loglevel            = SG_INFO;
  DataSocket int
  //TelnetPort int
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
  
  //ReInitTelnet bool
  
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

//--------------------------------------------------------------------------

// Consruct and return pointer to new FG_SERVER instance
func NewFG_SERVER() *FG_SERVER {
	ob := new(FG_SERVER)
	ob.ServerVersion = &Version{Major: 1, Minor: 1} // TODO
	ob.ProtocolVersion = &Version{Major: 1, Minor: 1} // TODO
	
	ob.BlackList = make(map[string]bool)
	ob.Telnet = NewTelnetServer()
	
	 	
	// set other defaults here
	return ob
}

//--------------------------------------------------------------------------

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
	me.Telnet.Port = port
	me.Telnet.Reinit = true
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



//////////////////////////////////////////////////////////////////////
/**
 * @brief Basic initialization
 * 
 *  If we are already initialized, close
 *  all connections and re-init all variables
 */

 
 
func (me *FG_SERVER) Init() error {
  //if LogFile != "" {
   // m_LogFile.open (m_LogFileName.c_str(), ios::out|ios::app);
    //sglog().setLogLevels( SG_ALL, SG_INFO );
   // sglog().enable_with_date (true);
   // sglog().set_output(m_LogFile);
  //}
  
  /*
  if (m_Initialized == true)
  {
    if (m_Listening)
    {
      Done();
    }
    m_Initialized       = false;
    m_Listening         = false;
    m_DataSocket        = 0;
    m_NumMaxClients     = 0;
    m_NumCurrentClients = 0;
  }
  */
  //if (m_ReinitData || m_ReinitTelnet)
  //{
  //  netInit ();
  //}
  /*
  if (m_ReinitData)
  {
    if (m_DataSocket)
    {
      delete m_DataSocket;
      m_DataSocket = 0;
    }
    m_DataSocket = new netSocket();
    if (m_DataSocket->open (false) == 0)    // UDP-Socket
    {
      SG_ALERT (SG_SYSTEMS, SG_ALERT, "FG_SERVER::Init() - "
        << "failed to create listener socket");
      return (ERROR_CREATE_SOCKET);
    }
    m_DataSocket->setBlocking (false);
    m_DataSocket->setSockOpt (SO_REUSEADDR, true);
    if (m_DataSocket->bind (m_BindAddress.c_str(), m_ListenPort) != 0)
    {
      SG_ALERT (SG_SYSTEMS, SG_ALERT, "FG_SERVER::Init() - "
        << "failed to bind to port " << m_ListenPort);
      SG_ALERT (SG_SYSTEMS, SG_ALERT, "already in use?");
      return (ERROR_COULDNT_BIND);
    }
    m_ReinitData = false;
  }
  */
  if me.Telnet.Reinit {
    //if (m_TelnetSocket)
    //{
     // delete m_TelnetSocket;
      //m_TelnetSocket = 0;
    //}
    //m_TelnetSocket = 0;
    if me.Telnet.Port != 0 {
    	s := fmt.Sprintf(":%s", me.Telnet.Port ) // TODO ip address = 0.0.0.0 ?
      	ln, err := net.Listen("tcp", s)
      
      //if (m_TelnetSocket->open (true) == 0)   // TCP-Socket
      //{
      //  SG_ALERT (SG_SYSTEMS, SG_ALERT, "FG_SERVER::Init() - "
      //    << "failed to create telnet socket");
      //  return (ERROR_CREATE_SOCKET);
     // }
     if err != nil {
     	log.Fatal("Cannot create telnet socket")
     	return err
     }
     for {
     	conn, err := ln.Accept() 
     	if err != nil {
     		log.Println(err)
     	}
     	go me.HandleTelnet(conn)
     }
      /*m_TelnetSocket->setBlocking (false);
      m_TelnetSocket->setSockOpt (SO_REUSEADDR, true);
      if (m_TelnetSocket->bind (m_BindAddress.c_str(), m_TelnetPort) != 0)
      {
        SG_ALERT (SG_SYSTEMS, SG_ALERT, "FG_SERVER::Init() - "
          << "failed to bind to port " << m_TelnetPort);
        SG_ALERT (SG_SYSTEMS, SG_ALERT, "already in use?");
        return (ERROR_COULDNT_BIND);
      }
      if (m_TelnetSocket->listen (MAX_TELNETS) != 0)
      {
        SG_ALERT (SG_SYSTEMS, SG_ALERT, "FG_SERVER::Init() - "
          << "failed to listen to telnet port");
        return (ERROR_COULDNT_LISTEN);
      }*/
    }
    me.Telnet.Reinit = false
  }
  
  /*
  SG_ALERT (SG_SYSTEMS, SG_ALERT, "# This is " << m_ServerName);
  SG_ALERT (SG_SYSTEMS, SG_ALERT, "# FlightGear Multiplayer Server v"
    << VERSION << " started");
  SG_ALERT (SG_SYSTEMS, SG_ALERT, "# using protocol version v"
    << m_ProtoMajorVersion << "." << m_ProtoMinorVersion
    << " (LazyRelay enabled)");
  if (m_BindAddress != "")
  {
    SG_ALERT (SG_SYSTEMS, SG_ALERT,"# listening on " << m_BindAddress);
  }
  SG_ALERT (SG_SYSTEMS, SG_ALERT,"# listening to port " << m_ListenPort);
  SG_ALERT (SG_SYSTEMS, SG_ALERT,"# telnet port " << m_TelnetPort);
  SG_ALERT (SG_SYSTEMS, SG_ALERT,"# using logfile " << m_LogFileName);
  if (m_IamHUB)
  {
    SG_ALERT (SG_SYSTEMS, SG_ALERT, "# I am a HUB Server");
  }
  */
  /*
  if (m_IsTracked)
  {
    if ( m_Tracker->InitTracker(&m_TrackerPID) )
    {
        SG_ALERT (SG_SYSTEMS, SG_ALERT, "# InitTracker FAILED! Disabling tracker!");
            m_IsTracked = false;
    }
    else
    {
#ifdef USE_TRACKER_PORT
      SG_ALERT (SG_SYSTEMS, SG_ALERT, "# tracked to "
        << m_Tracker->GetTrackerServer ()
        << ":" << m_Tracker->GetTrackerPort ()
        << ", using a thread." );
#else // #ifdef USE_TRACKER_PORT
      SG_ALERT (SG_SYSTEMS, SG_ALERT, "# tracked to "
        << m_Tracker->GetTrackerServer ()
        << ":" << m_Tracker->GetTrackerPort ());
#endif // #ifdef USE_TRACKER_PORT y/n
    }
  }
  else
  {
    SG_ALERT (SG_SYSTEMS, SG_ALERT, "# tracking is disabled.");
  }
  */
  /*
  SG_ALERT (SG_SYSTEMS, SG_ALERT, "# I have " << m_RelayList.size() << " relays");
  mT_RelayListIt CurrentRelay = m_RelayList.begin();
  while (CurrentRelay != m_RelayList.end())
  {
    SG_ALERT (SG_SYSTEMS, SG_ALERT, "# relay " << CurrentRelay->Name);
    CurrentRelay++;
  }
  SG_ALERT (SG_SYSTEMS, SG_ALERT, "# I have " << m_CrossfeedList.size() << " crossfeeds");
  mT_RelayListIt CurrentCrossfeed = m_CrossfeedList.begin();
  while (CurrentCrossfeed != m_CrossfeedList.end())
  {
    SG_ALERT (SG_SYSTEMS, SG_ALERT, "# crossfeed " << CurrentCrossfeed->Name
      << ":" << CurrentCrossfeed->Address.getPort());
    CurrentCrossfeed++;
  }
  SG_ALERT (SG_SYSTEMS, SG_ALERT, "# I have " << m_BlackList.size() << " blacklisted IPs");
  SG_ALERT (SG_SYSTEMS, SG_ALERT, "# Files: exit=[" << exit_file << "] stat=[" << stat_file << "]");
  m_Listening = true;
  return (SUCCESS);
  */
	return nil
} // FG_SERVER::Init()

/**
 *  Handle a telnet session. if a telnet connection is opened, this 
 *  method outputs a list  of all known clients.
 */
func (me *FG_SERVER) HandleTelnet(conn net.Conn){

	//var errno int = 0
	var Message string  = ""
  
  	
  /** @brief  Geodetic Coordinates */
  //Point3D         PlayerPosGeod;  
  //FG_Player CurrentPlayer;
  //netSocket       NewTelnet;
  //unsigned int  it;
  //NewTelnet.setHandle (Fd);
  //errno = 0;
  //////////////////////////////////////////////////
  //
  //      create the output message
  //      header
  //
  //////////////////////////////////////////////////
  Message  = "# This is " + me.ServerName
  Message += "\n"
  //Message += "# FlightGear Multiplayer Server v" + string(Me.VERSION);
  Message += "\n"
  Message += "# using protocol version v"
  //Message += NumToStr (me.ProtoMajorVersion, 0)
  //Message += "." + NumToStr (me.ProtoMinorVersion, 0)
  Message += " (LazyRelay enabled)"
  Message += "\n"
  /* if ( m_IsTracked )
  {
    Message += "# This server is tracked: ";
    Message += m_Tracker->GetTrackerServer();
    Message += "\n";
  }
  if (NewTelnet.send (Message.c_str(),Message.size(), MSG_NOSIGNAL) < 0)
  {
    if ((errno != EAGAIN) && (errno != EPIPE))
    {
      SG_LOG (SG_SYSTEMS, SG_ALERT, "FG_SERVER::HandleTelnet() - " << strerror (errno));
    }
    return (0);
  } */
  /* pthread_mutex_lock (& m_PlayerMutex);
  Message  = "# "+ NumToStr (m_PlayerList.size(), 0);
  pthread_mutex_unlock (& m_PlayerMutex);
  Message += " pilot(s) online\n";
  if (NewTelnet.send (Message.c_str(),Message.size(), MSG_NOSIGNAL) < 0)
  {
    if ((errno != EAGAIN) && (errno != EPIPE))
    {
      SG_LOG (SG_SYSTEMS, SG_ALERT, "FG_SERVER::HandleTelnet() - " << strerror (errno));
    }
    return (0);
  }*/
  //////////////////////////////////////////////////
  //
  //      create list of players
  //
  //////////////////////////////////////////////////
  /*
  it = 0;
  for (;;)
  {
    pthread_mutex_lock (& m_PlayerMutex);
    if (it < m_PlayerList.size())
    {
      CurrentPlayer = m_PlayerList[it]; 
      it++;
    }
    else
    {
      pthread_mutex_unlock (& m_PlayerMutex);
      break;
    }
    pthread_mutex_unlock (& m_PlayerMutex);
    sgCartToGeod (CurrentPlayer.LastPos, PlayerPosGeod);
    Message = CurrentPlayer.Callsign + "@";
    if (CurrentPlayer.IsLocal)
    {
      Message += "LOCAL: ";
    }
    else
    {
      mT_RelayMapIt Relay = m_RelayMap.find(CurrentPlayer.Address.getIP());
      if (Relay != m_RelayMap.end())
      {
        Message += Relay->second + ": ";
      }
      else
      {
        Message += CurrentPlayer.Origin + ": ";
      }
    }
    if (CurrentPlayer.Error != "")
    {
      Message += CurrentPlayer.Error + " ";
    }
    Message += NumToStr (CurrentPlayer.LastPos[X], 6)+" ";
    Message += NumToStr (CurrentPlayer.LastPos[Y], 6)+" ";
    Message += NumToStr (CurrentPlayer.LastPos[Z], 6)+" ";
    Message += NumToStr (PlayerPosGeod[Lat], 6)+" ";
    Message += NumToStr (PlayerPosGeod[Lon], 6)+" ";
    Message += NumToStr (PlayerPosGeod[Alt], 6)+" ";
    Message += NumToStr (CurrentPlayer.LastOrientation[X], 6)+" ";
    Message += NumToStr (CurrentPlayer.LastOrientation[Y], 6)+" ";
    Message += NumToStr (CurrentPlayer.LastOrientation[Z], 6)+" ";
    Message += CurrentPlayer.ModelName;
    Message += "\n";
    if (NewTelnet.send (Message.c_str(),Message.size(), MSG_NOSIGNAL) < 0)
    {
      if ((errno != EAGAIN) && (errno != EPIPE))
      {
        SG_LOG (SG_SYSTEMS, SG_ALERT, "FG_SERVER::HandleTelnet() - " << strerror (errno));
      }
      return (0);
    }
  }*/
 // NewTelnet.close ();
  //return (0);
} // FG_SERVER::HandleTelnet ()

