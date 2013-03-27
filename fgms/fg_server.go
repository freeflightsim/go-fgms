
package fgms

import(
	//"bytes"
	"fmt"
	"log"
	"net"		
	//"strings"
	//"strconv"
	//"time"
	//"unsafe"
)

import(
	//"github.com/davecgh/go-xdr/xdr"

	"github.com/fgx/go-fgms/tracker"
	"github.com/fgx/go-fgms/flightgear"
)

const VERSION = "0.1-go-experimental"

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
const RELAY_MAGIC             = 0x53464746    // GSGF


const (
	SENDER_UNKNOWN  = 0
	SENDER_KNOWN   
	SENDER_DIFF_IP  // Not sure this is used
)
		
// Main Server
type FG_SERVER struct {

	/*typedef union
	{
		uint32_t    complete;
		int16_t     High;
		int16_t     Low;
	} converter; 
	converter*    tmp; 
	*/

	//ServerVersion *Version
	
	Initialized bool
	
	ReinitData bool

	Listening bool


	ServerName  string
	BindAddress string
	ListenPort int
	IamHUB bool
	

	//PlayerList []*FG_Player
	Players map[string]*FG_Player
	//PlayerList map[string]*FG_Player
	PlayerExpires int
	
	

	Telnet *TelnetServer
	DataSocket *net.UDPConn

	//Loglevel            = SG_INFO; / TODO
	LogFileName string

	NumMaxClients int
	PlayerIsOutOfReach int // nautical miles
	NumCurrentClients int
	IsParent bool
	MaxClientID int


	
	//tmp                   = (converter*) (& PROTO_VER);
	//ProtoMinorVersion   = tmp->High;
	//ProtoMajorVersion   = tmp->Low;
	//LogFileName         = DEF_SERVER_LOG; // "fg_server.log";
	//wp                  = fopen("wp.txt", "w");

	//= maybe this could be a slice
	BlackList map[string]bool
	BlackListRejected uint64

	//RelayList map[string]*net.UDPConn
	//RelayMap map[string]*net.UDPConn
	//RelayList []*NetAddress
	Relays map[string]*net.UDPConn
	Crossfeeds map[string]*net.UDPConn

	IsTracked bool
	Tracker *tracker.FG_TRACKER

	MessageList []flightgear.T_ChatMsg

	//UpdateSecs          = DEF_UPDATE_SECS;
	// clear stats - should show what type of packet was received
	PacketsReceived int
	TelnetReceived int 
	BlackRejected int
	PacketsInvalid int //     = 0;  // invalid packet
	PktsForwarded int
	
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



// Construct and return pointer to new FG_SERVER instance
func NewFG_SERVER() *FG_SERVER {
	ob := new(FG_SERVER)
	
	ob.Players = make(map[string]*FG_Player)
	//ob.PlayerList = make([]*FG_Player, 0)
		
	//ob.RelayList = make([]*NetAddress, 0)
	//ob.RelayMap = make(map[string]string)
	ob.Relays = make(map[string]*net.UDPConn, 0)
	ob.Crossfeeds = make(map[string]*net.UDPConn, 0)
	
	ob.BlackList = make(map[string]bool)
		
	ob.Telnet = NewTelnetServer()
	
	ob.MessageList = make([]flightgear.T_ChatMsg, 0)
		
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
	log.Println("> SetDataPort=", port)
	me.ListenPort = port
	me.ReinitData = true
}

func (me *FG_SERVER) SetTelnetPort(port int){
	log.Println("> SetTelnetPort=", port)
	me.Telnet.Port = port
	me.Telnet.Reinit = true
}

// Set nautical miles two players must be apart to be out of reach
func (me *FG_SERVER) SetOutOfReach(nm int){
	log.Println("> SetOutOfReach=", nm, " nm")
	me.PlayerIsOutOfReach = nm
}

// Set time in seconds. if no packet arrives from a client
// within this time, the connection is dropped.  
func (me *FG_SERVER) SetPlayerExpires(secs int){
	log.Println("> SetPlayerExpires=", secs, " secs")
	me.PlayerExpires = secs
}

// Set if we are running as a Hubserver
func (me *FG_SERVER) SetHub(am_hub bool){
	log.Println("> SetHub=", am_hub)
	me.IamHUB = am_hub
}

// Set the logfile - TODO LOg FIle writing etc
func (me *FG_SERVER) SetLogfile( log_file_name string){
	log.Println("> SetLogfile=", log_file_name)
	me.LogFileName = log_file_name
	//TODO after research of simgear
}







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





// ---------------------------------------------------------------------------

// Basic initialization. 
// - TODO: If we are already initialized, close
// all connections and re-init all variables
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

	if me.ReinitData {
		if me.DataSocket != nil {
			//delete m_DataSocket;
			//m_DataSocket = 0;
		}
				
		//=== UDP ===
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5000")
		var erru error
		me.DataSocket, erru = net.ListenUDP("udp", addr)
		if erru != nil {
			log.Panicf("Fatal error starting UDP server: %s", erru)
			return err
		}
		me.ReinitData = false
	}

	if me.Telnet.Reinit {
	
		if me.Telnet.Port > 0 {
			me.Telnet.Addr = fmt.Sprintf(":%d", me.Telnet.Port ) // TODO ip address = 0.0.0.0 ?
			//= Create and listen on Telnet Socket
			var err error
			me.Telnet.Listen, err = net.Listen("tcp", me.Telnet.Addr)
			if err != nil {
				log.Panicf("Error Opening Telnet: %s", err)
			}
			me.Telnet.Reinit = false
		}
	}
	
	
	//log.Fatal("HERE")
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







//////////////////////////////////////////////////////////////////////
/**
* @brief  If we receive bad data from a client, we add the client to
*         the internal list anyway, but mark them as bad. But first 
*          we look if it isn't already there.
*          Send an error message to the bad client.
* @param Sender
* @param ErrorMsg
* @param IsLocal
*/
func (me *FG_SERVER) AddBadClient(Sender *net.UDPAddr , ErrorMsg string, IsLocal bool){
	//TODO
	//string                  Message;
	//FG_Player               NewPlayer;
	//mT_PlayerListIt         CurrentPlayer;
	
	//CurrentPlayer = m_PlayerList.begin();
	//////////////////////////////////////////////////
	//      see, if we already know the client
	//////////////////////////////////////////////////
	/* while (CurrentPlayer != m_PlayerList.end())
	{
		if (CurrentPlayer->Address.getIP() == Sender.getIP())
		{
		CurrentPlayer->Timestamp = time (0);
		return;
		}
		CurrentPlayer++;
	} */
	//////////////////////////////////////////////////
	//      new client, send an error message
	//////////////////////////////////////////////////
	fmt.Println("BADCLIENT", Sender)
	me.MaxClientID++
	NewPlayer := NewFG_Player()
	NewPlayer.Callsign      = "* Bad Client *"
	NewPlayer.ModelName     = "* unknown *"
	//NewPlayer.Timestamp     = time(0);
	//NewPlayer.JoinTime      = NewPlayer.Timestamp;
	// NewPlayer.Origin        = Sender.Host //getHost ()
	//NewPlayer.Address       = Sender.Address
	NewPlayer.IsLocal       = IsLocal
	NewPlayer.HasErrors     = true
	NewPlayer.Error         = ErrorMsg
	NewPlayer.ClientID      = me.MaxClientID
	//NewPlayer.PktsReceivedFrom      = 0
	//NewPlayer.PktsSentTo            = 0
	//NewPlayer.PktsForwarded         = 0
	//NewPlayer.LastRelayedToInactive = 0
	//SG_LOG (SG_SYSTEMS, SG_ALERT, "FG_SERVER::AddBadClient() - " << ErrorMsg);
	//Message = "bad client connected: ";
	//Message += Sender.getHost() + string(": ");
	//Message += ErrorMsg;
	//CreateChatMessage (NewPlayer.ClientID, Message);
	//pthread_mutex_lock (& m_PlayerMutex);
	//m_PlayerList.push_back (NewPlayer);
	//m_NumCurrentClients++;
	//pthread_mutex_unlock (& m_PlayerMutex);
	//*/
	//me.PlayerList[Sender.Ip] = NewPlayer
} // FG_SERVER::AddBadClient ()









// Main Loop
func (me *FG_SERVER) Loop() {

	//== Startup Telnet Listener
	go func(lisTel net.Listener){
		for {
			conna, erra := lisTel.Accept() 
			if erra != nil {
				log.Println(erra)
			}
			go me.HandleTelnetData(conna)
		}
	}(me.Telnet.Listen)
	log.Println("# Listening Telnet > ")
	
	
	//== Startup UDP listener
	count := 0
	buf := make([]byte, MAX_PACKET_SIZE)
	log.Println("# Listening UDP > ", )
	for {
		length, raddr, err := me.DataSocket.ReadFromUDP(buf)
		if err != nil {
				log.Printf("ReadFrom: %v", err)
				//break
		}else {
			count++
			//log.Printf("<%s> %q", raddr, buf[:length])
			//log.Println("count", count, raddr, length)
			//log.Println(buf[:length])
			//Msg []byte, Bytes int, SenderAddress *NetAddress){
			me.HandlePacket( buf[:length], length, raddr)
			
		}
	}
	fmt.Println("Should Never See This")
		
}


//////////////////////////////////////////////////////////////////////
// Look if we know the sending client
// return - 0: Sender is unknown  - 1: Sender is known - 2: Sender is known, but has a different IP
func (me *FG_SERVER) SenderIsKnown(senderCallsign string) int {

	//addr := SenderAddress.String()
	fmt.Println("Find=", senderCallsign)
	_, found := me.Players[senderCallsign]
	//for _, player := range me.PlayerList {
	//	if player.Callsign == SenderCallsign {
	//		//if player.
	//	}
	//	
	//}
	if found {
		return SENDER_KNOWN
	}
	/* mT_PlayerListIt CurrentPlayer;
	for (CurrentPlayer = m_PlayerList.begin();
	CurrentPlayer != m_PlayerList.end();
	CurrentPlayer++)
	{
	if (CurrentPlayer->Callsign == SenderCallsign){
		if CurrentPlayer->Address.getIP() == SenderAddress.getIP() {
			return 1 // Sender is known
		}
		// Same callsign, but different IP.
		// Quietly ignore this packet.
		return 2
		}
	} */
	// Sender is unkown
	return SENDER_UNKNOWN
} // FG_SERVER::SenderIsKnown ()







//////////////////////////////////////////////////////////////////////
//  Insert a new client to internal list
func (me *FG_SERVER) AddClient(Sender *net.UDPAddr, MsgHdr flightgear.T_MsgHdr, PosMsg flightgear.T_PositionMsg) {
	//time_t          Timestamp;
	//uint32_t        MsgLen;
	//uint32_t        MsgId;
	//uint32_t        MsgMagic;
	//string          Message;
	//string          Origin;
	//T_MsgHdr*       MsgHdr;
	//T_PositionMsg*  PosMsg;
	//FG_Player       NewPlayer;
	//bool    IsLocal;

	//Timestamp           = time(0);
	//MsgHdr              = (T_MsgHdr *) Msg;
	//var MsgHdr &flightgear.T_MsgTdr{}
	//PosMsg              = (T_PositionMsg *) (Msg + sizeof(T_MsgHdr));
	//MsgId               = XDR_decode<uint32_t> (MsgHdr->MsgId);
	//MsgLen              = XDR_decode<uint32_t> (MsgHdr->MsgLen);
	//MsgMagic            = XDR_decode<uint32_t> (MsgHdr->Magic);
	//IsLocal             = true;
	
	
	IsLocal := MsgHdr.Magic != RELAY_MAGIC  // not a local client
	fmt.Println (" ADD Client", Sender, IsLocal,  len(me.Players))
		
	var callsign string = MsgHdr.CallsignString()
	NewPlayer := NewFG_Player()
	NewPlayer.Callsign  = callsign
	NewPlayer.Passwd    = "test" //MsgHdr->Passwd;
	NewPlayer.ModelName = PosMsg.ModelString()
	//NewPlayer.Timestamp = time.Now().Unix()
	//NewPlayer.JoinTime  = NewPlayer.Timestamp
	//NewPlayer.Origin    = Sender.getHost () TODO
	NewPlayer.HasErrors = false
	// NewPlayer.Address   = Sender
	NewPlayer.IsLocal   = IsLocal
	//NewPlayer.LastPos.Clear()
	//NewPlayer.LastOrientation.Clear()
	//NewPlayer.PktsReceivedFrom = 0
	//NewPlayer.PktsSentTo       = 0
	//NewPlayer.PktsForwarded    = 0
	//NewPlayer.LastRelayedToInactive = 0 
	/* NewPlayer.LastPos.Set (
		XDR_decode64<double> (PosMsg->position[X]),
		XDR_decode64<double> (PosMsg->position[Y]),
		XDR_decode64<double> (PosMsg->position[Z])
	); */
	/*NewPlayer.LastOrientation.Set (
		XDR_decode<float> (PosMsg->orientation[X]),
		XDR_decode<float> (PosMsg->orientation[Y]),
		XDR_decode<float> (PosMsg->orientation[Z])
	);*/
	//NewPlayer.ModelName = PosMsg.ModelString()
	//m_MaxClientID++
	NewPlayer.ClientID = me.MaxClientID
	//pthread_mutex_lock (& m_PlayerMutex)
	//m_PlayerList.push_back (NewPlayer)
	
	
	//me.PlayerList = append(me.PlayerList, NewPlayer)
	//pthread_mutex_unlock (& m_PlayerMutex);
	
	// Add to List
	me.Players[callsign] = NewPlayer
	me.NumCurrentClients++
	if me.NumCurrentClients > me.NumMaxClients {
		me.NumMaxClients = me.NumCurrentClients;
	}
	var Message string = ""
	if IsLocal {
		Message  = "Welcome to "
		Message += me.ServerName
		me.CreateChatMessage (NewPlayer.ClientID , Message)
		//Message = "this is version v" + string(VERSION)
		//Message += " (LazyRelay enabled)"
		//CreateChatMessage (NewPlayer.ClientID , Message)
		//Message  ="using protocol version v"
		//Message += NumToStr (m_ProtoMajorVersion, 0)
		//Message += "." + NumToStr (m_ProtoMinorVersion, 0)
		//if me.IsTracked {
		//	Message += "This server is tracked."
		//}
		me.CreateChatMessage (NewPlayer.ClientID , Message)
		//UpdateTracker (NewPlayer.Callsign, NewPlayer.Passwd,
		//NewPlayer.ModelName, NewPlayer.Timestamp, CONNECT); 
		
	}
	/* Message  = NewPlayer.Callsign;
	Message += " is now online, using ";
	CreateChatMessage (0, Message);
	Message  = NewPlayer.ModelName;
	CreateChatMessage (0, Message);
	Origin  = NewPlayer.Origin;
	if IsLocal{
		Message = "New LOCAL Client: ";
	}else{
		Message = "New REMOTE Client: ";
		mT_RelayMapIt Relay = m_RelayMap.find(NewPlayer.Address.getIP());
		if (Relay != m_RelayMap.end())
		{
		Origin = Relay->second;
		}
	} */
	/*
	SG_LOG (SG_SYSTEMS, SG_INFO, Message
		<< NewPlayer.Callsign << " "
		<< Origin << ":" << Sender.getPort()
		<< " (" << NewPlayer.ModelName << ")"
		<< " current clients: "
		<< m_NumCurrentClients << " max: " << m_NumMaxClients
	); */
} // FG_SERVER::AddClient()
