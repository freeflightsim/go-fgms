
package fgms

import(
	//"bytes"
	"fmt"
	"log"
	"net"
	"path/filepath"
	//"strings"
	//"strconv"
	//"time"
	//"unsafe"

	"github.com/freeflightsim/go-fgms/tracker"
	"github.com/freeflightsim/go-fgms/message"

)


		
// Main Server
type FgServer struct {


	//ServerVersion *Version
	
	Initialized bool `json:"initialised"`
	
	ReinitData bool

	Listening bool

	ServerName  string `json:"server_name"`
	BindAddress string `json:"address"`
	ListenPort int `json:"port"`
	IamHUB bool `json:"is_hub"`
	

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
	//LogFileName         = DEF_SERVER_LOG; // "FgServer.log";
	//wp                  = fopen("wp.txt", "w");

	//= maybe this could be a slice
	//BlackList *blacklist
	//BlackListRejected uint64

	//RelayList map[string]*net.UDPConn
	//RelayMap map[string]*net.UDPConn
	//RelayList []*NetAddress
	Relays map[string]*net.UDPConn
	
	
	//Crossfeeds map[string]*UDP_Conn //*net.UDPConn
	//CrossFeedFailed int
	//CrossFeedSent int
	//MT_CrossFeedFailed int
	//MT_CrossFeedSent int

	
	IsTracked bool
	Tracker *tracker.FG_Tracker

	MessageList []message.ChatMsg

	//UpdateSecs          = DEF_UPDATE_SECS;
	// clear stats - should show what type of packet was received
	PacketsReceived int
	PacketsInvalid int //     = 0;  // invalid packet
	
	TelnetReceived int 
	//BlackRejected int
	
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


	TrackerConnect int
	TrackerDisconnect int
	TrackerPostion int // Tracker messages queued
	//pthread_mutex_init( &m_PlayerMutex, 0 );
} 

var Server *FgServer

func SetupServer(){


	//fmt.Println("AUTO SERVER")
	Server = new(FgServer)

	//Server.BlackList = blacklist{}
	//Server.BlackList.Hosts = make(map[string]bool, 0)

	Server.Players = make(map[string]*FG_Player)
	//ob.PlayerList = make([]*FG_Player, 0)
		
	//ob.RelayList = make([]*NetAddress, 0)
	//ob.RelayMap = make(map[string]string)
	Server.Relays = make(map[string]*net.UDPConn, 0)
	//ob.Crossfeeds = make(map[string]*net.UDPConn, 0)
	//ob.Crossfeeds = make(map[string]*UDP_Conn, 0)
	
	//ob.BlackList = make(map[string]bool)

	Server.Telnet = NewTelnetServer()

	Server.MessageList = make([]message.ChatMsg, 0)

	//Server.Init()

	//InitBlacklist()
	//InitCrossfeed()
	//InitHttp()
}


 

//--------------------------------------------------------------------------

func (me *FgServer) SetServerName(name string){
	me.ServerName = name
}
func (me *FgServer) SetBindAddress(addr string){
	me.BindAddress = addr
}

func (me *FgServer) SetDataPort(port int){
	me.ListenPort = port
	me.ReinitData = true
}

func (me *FgServer) SetTelnetPort(port int){
	me.Telnet.Port = port
	me.Telnet.Reinit = true
}

// Set nautical miles two players must be apart to be out of reach
func (me *FgServer) SetOutOfReach(nm int){
	me.PlayerIsOutOfReach = nm
}

// Set time in seconds. if no packet arrives from a client
// within this time, the connection is dropped.  
func (me *FgServer) SetPlayerExpires(secs int){
	me.PlayerExpires = secs
}

// Set if we are running as a Hubserver
func (me *FgServer) SetHub(am_hub bool){
	me.IamHUB = am_hub
}

// Set the logfile - TODO LOg FIle writing etc
func (me *FgServer) SetLogfile( log_file_name string){
	//log.Println("> SetLogfile=", log_file_name)
	me.LogFileName = log_file_name
	//TODO after research of simgear
}











// ---------------------------------------------------------------------------

// Basic initialization. 
// - TODO: If we are already initialized, close
// all connections and re-init all variables
func (me *FgServer) Init() error {
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
	fmt.Println("Server.INIT()")
	if me.ReinitData {
		if me.DataSocket != nil {
			//delete m_DataSocket;
			//m_DataSocket = 0;
		}
				
		//=== UDP ===
		addr, err := net.ResolveUDPAddr("udp", "192.168.50.5:5000")
		var erru error
		me.DataSocket, erru = net.ListenUDP("udp", addr)
		if erru != nil {
			log.Panicf("Fatal error starting UDP server: %s", erru)
			return err
		}
		fmt.Println("SOCKET INIT------------------------------")
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
} // FgServer::Init()


// Read a config file and set internal variables accordingly.

func SetConfig(conf Config){
	Server.SetConfig(conf)
}

func (me *FgServer) SetConfig(conf Config) error {

	fmt.Println("Server.SetConfig()", conf.Server, me)
	// Server Name
	me.SetServerName(conf.Server.Name)

	// Address
	me.SetBindAddress(conf.Server.Address)

	// UDP Port No
	me.SetDataPort(conf.Server.Port)

	// Telnet Port
	me.SetTelnetPort(conf.Server.TelnetPort)

	// Outta Reach
	me.SetOutOfReach(conf.Server.OutOfReachNm)

	// Player Expires
	me.SetPlayerExpires(conf.Server.PlayerExpiresSecs)

	// Server is hub
	me.SetHub( conf.Server.IsHub )

	// Log File
	me.SetLogfile(conf.Server.LogFile);

	// Tracked
	/*
	Val, err = conf.Get ("server.tracked")
	if Val != "" {
		tracked, _ := strconv.ParseBool(Val)
		if tracked {
			trkServer, err := conf.Get("server.tracking_server")
			if err != nil {
				log.Fatalln("Error", "Missing `server.tracking_server`", trkServer)
				return err
			}
			fmt.Println("TRK", trkServer,  tracked)
			me.AddTracker(trkServer, pii, tracked)

		}
	}
	*/
	//if true == true {
	//	return nil
	//}

	// Read the list of relays
	for _, relay := range conf.Relays {
		//me.AddRelay(relay.Host, relay.Port)
		fmt.Println(relay)
	}

	// Read the list of crossfeeds
	for _, cf := range conf.Crossfeeds {
		CrossFeed.Add(cf.Host, cf.Port)
	}


	 // read the list of blacklisted IPs
	for _, bl := range conf.Blacklists {
		Blacklist.Add(bl)
	}

	return nil
}






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
func (me *FgServer) AddBadClient(Sender *net.UDPAddr , ErrorMsg string, IsLocal bool){
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
	NewPlayer := new(FG_Player)
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
	//SG_LOG (SG_SYSTEMS, SG_ALERT, "FgServer::AddBadClient() - " << ErrorMsg);
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
} // FgServer::AddBadClient ()







// Main Loop
func (me *FgServer) Start() {

	//== Startup Telnet Listener
	/*
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
	*/
	
	
	//== Start loop to check ocassinally crossfeed server conentions
	//go me.StartCrossfeedCheckTimer()
	me.Init()
	
	//== Startup UDP listener
	count := 0
	buffer := make([]byte, MAX_PACKET_SIZE)
	log.Println("# Listening UDP > ", me.DataSocket)
	for {
		length, raddr, err := me.DataSocket.ReadFromUDP(buffer)
		if err != nil {
				log.Printf("ReadFrom: %v", err)
				//break
		}else {
			count++
			//log.Printf("<%s> %q", raddr, buf[:length])
			//log.Println("count", count, raddr, length)
			//log.Println(buf[:length])
			//Msg []byte, Bytes int, SenderAddress *NetAddress){
			me.HandlePacket( buffer[:length], length, raddr)
			
		}
	}
	fmt.Println("Should Never See This")
		
}


//////////////////////////////////////////////////////////////////////
// Look if we know the sending client
// return - 0: Sender is unknown  - 1: Sender is known - 2: Sender is known, but has a different IP
/*
func (me *FgServer) DEADSenderIsKnown(header message.HeaderMsg, address *net.UDPAddr) int {

	player, found := me.Players[header.Callsign()]

	if found == false {
		return SENDER_UNKNOWN
	}
	if player.Address.String() == address.String() {
		return SENDER_KNOWN
	}
	return SENDER_DIFF_IP
} // FgServer::SenderIsKnown ()
*/






//////////////////////////////////////////////////////////////////////
//  Insert a new client to internal list
//func (me *FgServer) AddClient(Sender *net.UDPAddr, MsgHdr flightgear.T_MsgHdr, PosMsg flightgear.T_PositionMsg) {
func (me *FgServer) AddClient(header *message.HeaderMsg, position *message.PositionMsg, address *net.UDPAddr, ) *FG_Player {


	var callsign string = header.Callsign()

	client := new(FG_Player)

	me.MaxClientID++
	client.ClientID = me.MaxClientID
	client.Address = address
	client.IsLocal = header.Magic != message.RELAY_MAGIC

	client.Timestamp = Now()
	client.JoinTime  = client.Timestamp

	client.Callsign  = callsign
	//client.Passwd    = "test" //MsgHdr->Passwd;

	//NewPlayer.Origin    = Sender.getHost () TODO

	client.ModelName = position.Model()
	s := filepath.Base(client.ModelName)
	client.Aircraft = s[0:len(s)-len(filepath.Ext(s))]




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
	//client.LastPos.Set( position.Position[X], position.Position[Y], position.Position[Z])
	 
	/*NewPlayer.LastOrientation.Set (
		XDR_decode<float> (PosMsg->orientation[X]),
		XDR_decode<float> (PosMsg->orientation[Y]),
		XDR_decode<float> (PosMsg->orientation[Z])
	);*/
	//client.LastOrientation.Set( float64(position.Orientation[X]), float64(position.Orientation[Y]), float64(position.Orientation[Z]))
	
	//NewPlayer.ModelName = PosMsg.ModelString()

	//pthread_mutex_lock (& m_PlayerMutex)
	//m_PlayerList.push_back (NewPlayer)
	//pthread_mutex_unlock (& m_PlayerMutex);
	
	// Add to Map, and increment counters

	me.Players[callsign] = client
	me.NumCurrentClients++
	if me.NumCurrentClients > me.NumMaxClients {
		me.NumMaxClients = me.NumCurrentClients
	}
	
	var Message string = ""
	if 1 ==2 && client.IsLocal {
		Message  = "Welcome to "
		Message += me.ServerName
		//me.CreateChatMessage (NewPlayer.ClientID , Message)
		
		Message = "this is version v" + VERSION
		Message += " (LazyRelay enabled)"
		//me.CreateChatMessage (NewPlayer.ClientID , Message)
		
		Message  ="using protocol version v" + GetProtocolVersionString()
		//Message += NumToStr (m_ProtoMajorVersion, 0)
		//Message += "." + NumToStr (m_ProtoMinorVersion, 0)
		if me.IsTracked {
			Message += "This server is tracked."
		}
		me.CreateChatMessage (client.ClientID , Message)
		//me.UpdateTracker(NewPlayer.Callsign, NewPlayer.Passwd, NewPlayer.ModelName, NewPlayer.Timestamp, tracker.CONNECT)
		 me.UpdateTracker(client, tracker.CONNECT)
		
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
	log.Println (" ADD Client ", callsign, address.String(), client.IsLocal,  len(me.Players))
	return client
	/*
	SG_LOG (SG_SYSTEMS, SG_INFO, Message
		<< NewPlayer.Callsign << " "
		<< Origin << ":" << Sender.getPort()
		<< " (" << NewPlayer.ModelName << ")"
		<< " current clients: "
		<< m_NumCurrentClients << " max: " << m_NumMaxClients
	); */
} // FgServer::AddClient()
