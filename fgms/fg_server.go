
package fgms

import(
	"fmt"
	"log"
	"net"
	"bytes"
	"bufio"
	"io"
	"strings"
)
import(
	"github.com/davecgh/go-xdr/xdr"

	"github.com/fgx/go-fgms/tracker"
	"github.com/fgx/go-fgms/flightgear"
	
	
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
	VERSION int
	ServerVersion *Version
	ProtocolVersion *Version
	
	Initialized bool
	
	ReinitData bool
	//ReinitTelnet bool
	Listening bool


	ServerName  string
	BindAddress string
	ListenPort int
	IamHUB bool
	


	PlayerList map[string]*FG_Player
	PlayerExpires int

	Telnet *TelnetServer
	TelnetChan chan net.Conn
	
	TelnetAdmin *TelnetServer
	TelnetAdminChan chan net.Conn

//Loglevel            = SG_INFO;
DataSocket net.Conn
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

	//= maybe this could be a list
	BlackList map[string]bool
	BlackListRejected uint64

	RelayMap map[string]string//           = map<uint32_t, string>();
	RelayList []*NetAddress

//typedef std::list<mT_Relay>               mT_RelayList;
//	typedef mT_RelayList::iterator            mT_RelayListIt;
//typedef std::map<uint32_t,string>         mT_IP2RelayNames;	
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

// Consrtruct and return pointer to new FG_SERVER instance
func NewFG_SERVER() *FG_SERVER {
	ob := new(FG_SERVER)
	ob.ServerVersion = &Version{Major: 1, Minor: 1} // TODO
	ob.ProtocolVersion = &Version{Major: 1, Minor: 1} // TODO
	
	ob.PlayerList = make(map[string]*FG_Player)
	
	ob.BlackList = make(map[string]bool)
	
	ob.RelayList = make([]*NetAddress, 0)
	ob.RelayMap = make(map[string]string)
	
	ob.Telnet = NewTelnetServer()
	//ob.TelnetChan = msgchan := make(chan string)
	
	ob.TelnetAdmin = NewTelnetServer()
	
		
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


// Set the logfile
func (me *FG_SERVER) SetLogfile( log_file_name string){
	log.Println("> SetLogfile=", log_file_name)
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
	
	// First create a relay object, which is a NetAddress
	NewRelay := NewMT_Relay(server, port)
	log.Println("> Add Relay = ", server, NewRelay.IpAddress)
	
	//= Now go and check it exists as IP as a callback
	go func(addr *NetAddress){
		err := NewRelay.LookupIP()
		if err != nil{
			log.Fatalln("    < Relay FAIL < No IP address for Host ", addr.Host, addr.Port)
			return 
		}
		me.RelayMap[NewRelay.Host] = NewRelay.IpAddress
		log.Println("    < Relay Added < Lookup OK:  ", addr.Host, NewRelay.IpAddress)
	}(NewRelay)	
	
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


// --------------------------------------------------------

// Add an IP to the blacklist
func (me *FG_SERVER) AddBlacklist(FourDottedIP string) {
	//SG_ALERT (SG_SYSTEMS, SG_ALERT, "Adding to blacklist: " << FourDottedIP);
	//m_BlackList[netAddress(FourDottedIP.c_str(), 0).getIP()] = true; // TODO lookup ip ?
	log.Println("> Add Blacklist = ", FourDottedIP)
	go func(ip_str string){
		addrs, err := net.LookupHost(ip_str)
		//err := net.LookupHost(ip_str)
		if err != nil{
			log.Fatalln("    < Blacklist FAIL: No IP address for address = ", ip_str)
			return 
		}
		log.Println("    < Blacklist Added < Lookup OK: ", ip_str, addrs, ip_str == addrs[0])
		
		me.BlackList[ addrs[0] ] = true
	}(FourDottedIP)
	
	//log.Println("Added to blacklist: ", FourDottedIP)
		
} 

// Check if the user is black listed. true if blacklisted
 func (me *FG_SERVER) IsBlackListed(SenderAddress *NetAddress) bool {
	_, ok :=  me.BlackList[SenderAddress.IpAddress]
	if ok {
		return true
	}
	return false
} 

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

	if 1 == 0 { //me.ReinitData {
		if me.DataSocket != nil {
			//delete m_DataSocket;
			//m_DataSocket = 0;
		}
		//m_DataSocket = new netSocket()
		//if (m_DataSocket->open (false) == 0)    // UDP-Socket
		//{
		//SG_ALERT (SG_SYSTEMS, SG_ALERT, "FG_SERVER::Init() - "
		//	<< "failed to create listener socket");
		//return (ERROR_CREATE_SOCKET);
		//	}
		/* s := fmt.Sprintf(":%d", 5000 )
		udp_ln, err := net.ListenUDP(s)
		if err != nil {
			log.Fatal("Cannot create UDP socket")
			return err
		}
		for {
				conn, err := udp_ln.ListenPacket() 
				if err != nil {
					log.Println(err)
				}
				log.Println(conn)
				//go me.HandleUDP(conn)
		}*/
		//m_DataSocket->setBlocking (false);
		//m_DataSocket->setSockOpt (SO_REUSEADDR, true);
		//if (m_DataSocket->bind (m_BindAddress.c_str(), m_ListenPort) != 0)
		//{
		//SG_ALERT (SG_SYSTEMS, SG_ALERT, "FG_SERVER::Init() - "
		//	<< "failed to bind to port " << m_ListenPort);
		//SG_ALERT (SG_SYSTEMS, SG_ALERT, "already in use?");
		//return (ERROR_COULDNT_BIND);
		//}
		me.ReinitData = false
	}

if me.Telnet.Reinit {
	//if (m_TelnetSocket)
	//{
	// delete m_TelnetSocket;
	//m_TelnetSocket = 0;
	//}
	//m_TelnetSocket = 0;
	if me.Telnet.Port != 0 {
		//s := fmt.Sprintf(":%d", me.Telnet.Port ) // TODO ip address = 0.0.0.0 ?
		//ln, err := net.Listen("tcp", s)
		//if err != nil {
		//	log.Fatal("Cannot create telnet socket")
		//	return err
		//}
		//telnetDataChan := make(chan TelnetClient)
		
		//TODO Make UDP Socket
		lsu, erruu := net.ListenPacket("udp", ":5000")
		if erru != nil{
			log.Fatal("Cannot create UDP socket")
		}
		
		
		// admin
		sa := fmt.Sprintf(":%d", 5005 ) 
		lna, erra := net.Listen("tcp", sa)
		if erra != nil {
			log.Fatal("Cannot create telnet socket for admin", erra)
			return erra
		}
		ta_msgChan := make(chan TelnetClient)
		ta_addChan := make(chan TelnetClient)
		ta_rmChan := make(chan TelnetClient)
		
		go me.H_TelnetAdminMessages( ta_msgChan, ta_addChan, ta_rmChan )
		
		//msgchan := make(chan string)
		
		//go me.PrintMessages(msgchan)
		for {
			//_, err := ln.Accept() 
			//if err != nil {
			//	log.Println(err)
			//}
			conna, erra := lna.Accept() 
			if erra != nil {
				log.Println(erra)
			}
			
			//go me.HandleTelnetData(conn, telnetDataChan)
			go me.HandleAdminTelnet(conna,  ta_msgChan, ta_addChan, ta_rmChan )
			
			
		}
		
	}
	me.Telnet.Reinit = false
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
log.Fatal("HERE")
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


//---------------------------------------------------------------------------



//---------------------------------------------------------------------------

/**
*  Handle a telnet session. if a telnet connection is opened, this 
*  method outputs a list  of all known clients.
*/
func (me *FG_SERVER) HandleTelnetData(conn net.Conn, telnetDataChan <-chan TelnetClient){

	//var errno int = 0
	var Message string  = ""
	//buf := make([]byte, 4096)

	
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
Message += "# FlightGear Multiplayer Server version: " + me.ServerVersion.Str()
Message += "\n"
Message += "# using protocol version: "
Message += me.ProtocolVersion.Str()
Message += " (LazyRelay enabled)"
Message += "\n"
//buf.Add

	// print conn.RemoteAddr()
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
	var buffer bytes.Buffer
	buffer.WriteString( Message )
	_, err := conn.Write( buffer.Bytes() )
	if err != nil {
		log.Println("error", err)
	}
	conn.Close()
	//return (0);
} // FG_SERVER::HandleTelnet ()

func PromptNick(c net.Conn, bufc *bufio.Reader) string {
	io.WriteString(c, "Welcome to fgms admin\n")
	io.WriteString(c, "fgms needs password > ")
	nick, _, _ := bufc.ReadLine()
	return string(nick)
}

/**
*  Handle a telnet session. if a telnet connection is opened, this 
*  method outputs a list  of all known clients.
*/
func (me *FG_SERVER) HandleAdminTelnet(c net.Conn, msgchan chan TelnetClient, 
										addchan chan<- TelnetClient, rmchan chan<- TelnetClient){
	log.Println("HandleAdminTelnet", c)
	bufc := bufio.NewReader(c)
	
	defer c.Close()
	client := TelnetClient{
		conn:     c,
		nickname: PromptNick(c, bufc),
		ch:       make(chan string),
	}
	if strings.TrimSpace(client.nickname) != "s" {
		io.WriteString(c, "Invalid Password\n")
		return
	}

	// Register user
	addchan <- client
	defer func() {
		//msgchan <- client //fmt.Sprintf("User %s left the chat room.\n", client.nickname)
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmchan <- client
	}()
	io.WriteString(c, fmt.Sprintf("Welcome, %s!\n\n", client.nickname))
	//msgchan <- fmt.Sprintf("New user %s has joined the chat room.\n", client.nickname)

	io.WriteString(c, "show | help | set var = value\n > ")
	// I/O
	go client.ReadLinesInto(msgchan)
	client.WriteLinesFrom(msgchan)
}



func (me *FG_SERVER) H_TelnetAdminMessages(msgChan <-chan TelnetClient, addChan <-chan TelnetClient, rmChan <-chan TelnetClient) {

	//clients := make(map[net.Conn]chan<- string)    
    for {
    	select {
    	
    		case client := <-msgChan:
    			log.Println("New Telnet Data request", client)
    			fmt.Println(">>", client.message)
    			m := strings.TrimSpace(client.message)
    			if m == "help" {
    					io.WriteString(client.conn, "HELP\n\n")
    				
    			}else if m == "show"{
					io.WriteString(client.conn, "SHOW\n\n")
    			}
    			//for _, ch := range clients {
				//	go func(mch chan<- string) { mch <- "\033[1;33;40m" + msg + "\033[m" }(ch)
				//}
    		
    		case client := <-addChan:
    			//log.Println("New Telnet Admin request", client)
    			log.Printf("New client: %v\n", client.conn)
				//clients[client.conn] = client.ch
				
			case client := <-rmChan:
				log.Printf("Client disconnects: %v\n", client.conn)
				//delete(clients, client.conn)
				
    	}
    }
  
}




//////////////////////////////////////////////////////////////////////
/**
 * @brief Handle client connections
 * @param Msg
 * @param Bytes
 * @param SenderAddress
 */
func (me *FG_SERVER) HandlePacket(Msg []byte, Bytes int, SenderAddress *NetAddress){
  //T_MsgHdr*       MsgHdr;
  //T_PositionMsg*  PosMsg;
  //uint32_t        MsgId;
  //uint32_t        MsgMagic;
  //Timestamp time.Time
  //Point3D         SenderPosition;
  //Point3D         SenderOrientation;
  //Point3D         PlayerPosGeod;
  //mT_PlayerListIt CurrentPlayer;
  //mT_PlayerListIt SendingPlayer;
  //unsigned int    PktsForwarded = 0;

  //Timestamp = time.Now() //time(0);
  //MsgHdr    = (T_MsgHdr *) Msg;
  	//MsgHdr :=  
  	var MsgHdr flightgear.T_MsgHdr
	remainingBytes, err := xdr.Unmarshal(Msg, &MsgHdr)
	fmt.Println("got", remainingBytes, err)
  	//MsgMagic  = XDR_decode<uint32_t> (MsgHdr->Magic);
  	//MsgId     = XDR_decode<uint32_t> (MsgHdr->MsgId);
  	fmt.Println( MsgHdr.Magic, MsgHdr.MsgId )
  	
  	//------------------------------------------------------
  	// First of all, send packet to all crossfeed servers.
  	//SendToCrossfeed (Msg, Bytes, SenderAddress);
  	//me.SendToCrossfeed(Msg, Bytes, SenderAddress)
  
  
  	//------------------------------------------------------
	//=  Now do the local processing
	if me.IsBlackListed(SenderAddress) {
	    me.BlackListRejected++
    	return
  	}
  	/*  WHY ??? passed by value
  	if ! me.PacketIsValid(	Bytes, 
  							MsgHdr, 
  							SenderAddress) {
    	me.PacketsInvalid++
    	return
  	} */
  	
  /* if (MsgMagic == RELAY_MAGIC) // not a local client
  {
    if (! IsKnownRelay (SenderAddress))
    {
      m_UnknownRelay++;
      return;
    }
    else
    {
      m_RelayMagic++; // bump relay magic packet
    }
  } */
  	if MsgHdr.Magic == RELAY_MAGIC {
  		if me.IsKnownRelay(SenderAddress) {
  			return
  		}
  		me.RelayMagic++ // bump relay magic packet
  }
  
  //////////////////////////////////////////////////
  //
  //    Store senders position
  //
  //////////////////////////////////////////////////
  /* if (MsgId == POS_DATA_ID)
  {
    m_PositionData++;
    PosMsg = (T_PositionMsg *) (Msg + sizeof(T_MsgHdr));
    double x = XDR_decode64<double> (PosMsg->position[X]);
    double y = XDR_decode64<double> (PosMsg->position[Y]);
    double z = XDR_decode64<double> (PosMsg->position[Z]);
    if ( (x == 0.0) || (y == 0.0) || (z == 0.0) )
    { // ignore while position is not settled
      return;
    }
    SenderPosition.Set (x, y, z);
    SenderOrientation.Set (
      XDR_decode<float> (PosMsg->orientation[X]),
      XDR_decode<float> (PosMsg->orientation[Y]),
      XDR_decode<float> (PosMsg->orientation[Z])
    );
  }
  else
  {
    m_NotPosData++;
  } */
  //////////////////////////////////////////////////
  //
  //    Add Client to list if its not known
  //
  //////////////////////////////////////////////////
  /* int ClientInList = SenderIsKnown (MsgHdr->Callsign, SenderAddress);
  if (ClientInList == 0)
  { // unknown, add to the list
    if (MsgId != POS_DATA_ID)
    { // ignore clients until we have a valid position
      return;
    }
    AddClient (SenderAddress, Msg);
  }
  else if (ClientInList == 2)
  { // known, but different IP => ignore
    return;
  }*/
  //////////////////////////////////////////
  //
  //      send the packet to all clients.
  //      since we are walking through the list,
  //      we look for the sending client, too. if it
  //      is not already there, add it to the list
  //
  //////////////////////////////////////////////////
  /* MsgHdr->Magic = XDR_encode<uint32_t> (MSG_MAGIC);
  SendingPlayer = m_PlayerList.end();
  CurrentPlayer = m_PlayerList.begin();
  while (CurrentPlayer != m_PlayerList.end())
  { */
    //////////////////////////////////////////////////
    //
    //      ignore clients with errors
    //
    //////////////////////////////////////////////////
    //if (CurrentPlayer->HasErrors)
    //{
     // CurrentPlayer++;
      //continue;
    //}
    //////////////////////////////////////////////////
    //        Sender == CurrentPlayer?
    //////////////////////////////////////////////////
    //  FIXME: if Sender is a Relay,
    //         CurrentPlayer->Address will be
    //         address of Relay and not the client's!
    //         so use a clientID instead
    /* if (CurrentPlayer->Callsign == MsgHdr->Callsign)
    {
      if (MsgId == POS_DATA_ID)
      {
        CurrentPlayer->LastPos         = SenderPosition;
        CurrentPlayer->LastOrientation = SenderOrientation;
      }
      else
      {
        SenderPosition    = CurrentPlayer->LastPos;
        SenderOrientation = CurrentPlayer->LastOrientation;
      }
      SendingPlayer = CurrentPlayer;
      CurrentPlayer->Timestamp = Timestamp;
      CurrentPlayer->PktsReceivedFrom++;
      CurrentPlayer++;
      continue; // don't send packet back to sender
    }*/
    //////////////////////////////////////////////////
    //      do not send packets to clients if the
    //      origin is an observer, but do send
    //      chat messages anyway
    //      FIXME: MAGIC = SFGF!
    //////////////////////////////////////////////////
    /* if ((strncasecmp(MsgHdr->Callsign, "obs", 3) == 0)
    &&  (MsgId != CHAT_MSG_ID))
    {
      return;
    } */
    //////////////////////////////////////////////////
    //
    //      do not send packet to clients which
    //      are out of reach.
    //
    //////////////////////////////////////////////////
    /* if ((Distance (SenderPosition, CurrentPlayer->LastPos) > m_PlayerIsOutOfReach)
    &&  (CurrentPlayer->Callsign.compare (0, 3, "obs", 3) != 0))
    {
      CurrentPlayer++;
      continue;
    } */
    //////////////////////////////////////////////////
    //
    //  only send packet to local clients
    //
    //////////////////////////////////////////////////
    /* if (CurrentPlayer->IsLocal)
    {
      SendChatMessages (CurrentPlayer);
      m_DataSocket->sendto (Msg, Bytes, 0, &CurrentPlayer->Address);
      CurrentPlayer->PktsSentTo++;
      PktsForwarded++;
    }
    CurrentPlayer++;
  } */
  /* 
  if (SendingPlayer == m_PlayerList.end())
  { // player not yet in our list
    // should not happen, but test just in case
    SG_LOG (SG_SYSTEMS, SG_ALERT, "## BAD => "
      << MsgHdr->Callsign << ":" << SenderAddress.getHost()
      << " : " << SenderIsKnown (MsgHdr->Callsign, SenderAddress)
    );
    return;
  }
  DeleteMessageQueue ();
  SendToRelays (Msg, Bytes, SendingPlayer);
  */
} // FG_SERVER::HandlePacket ( char* sMsg[MAX_PACKET_SIZE] )


func (me *FG_SERVER) PacketIsValid(	Bytes int, 
									MsgHdr flightgear.T_MsgHdr, 
									SenderAddress *NetAddress ) bool{

  //uint32_t        MsgMagic;
  //uint32_t        MsgLen;
  //uint32_t        MsgId;
  //string          ErrorMsg;
  //string          Origin;
  //typedef union
  //{
  //  uint32_t    complete;
  //  int16_t     High;
  //  int16_t     Low;
  //} converter;
  /* TODO	
  Origin = SenderAddress.getHost();
  MsgMagic = XDR_decode<uint32_t> (MsgHdr->Magic);
  MsgId  = XDR_decode<uint32_t> (MsgHdr->MsgId);
  MsgLen = XDR_decode<uint32_t> (MsgHdr->MsgLen);
  if (Bytes < (int)sizeof(MsgHdr))
  {
    ErrorMsg  = SenderAddress.getHost();
    ErrorMsg += " packet size is too small!";
    AddBadClient (SenderAddress, ErrorMsg, true);
    return (false);
  }
  */
  
  //= Check magic
  /*
  if ((MsgMagic != MSG_MAGIC) && (MsgMagic != RELAY_MAGIC))
  {
    char m[5];
    memcpy (m, (char*) &MsgMagic, 4);
    m[4] = 0;
    ErrorMsg  = Origin;
    ErrorMsg += " BAD magic number: ";
    ErrorMsg += m;
    AddBadClient (SenderAddress, ErrorMsg, true);
    return (false);
  }
  */
  	if MsgHdr.Magic != flightgear.MSG_MAGIC && MsgHdr.Magic != RELAY_MAGIC {
 		
  
	}  
	/*
  if (XDR_decode<uint32_t> (MsgHdr->Version) != PROTO_VER)
  {
    MsgHdr->Version = XDR_decode<uint32_t> (MsgHdr->Version);
    ErrorMsg  = Origin;
    ErrorMsg += " BAD protocol version! Should be ";
    converter*    tmp;
    tmp = (converter*) (& PROTO_VER);
    ErrorMsg += NumToStr (tmp->High, 0);
    ErrorMsg += "." + NumToStr (tmp->Low, 0);
    ErrorMsg += " but is ";
    tmp = (converter*) (& MsgHdr->Version);
    ErrorMsg += NumToStr (tmp->Low, 0);
    ErrorMsg += "." + NumToStr (tmp->High, 0);
    AddBadClient (SenderAddress, ErrorMsg, true);
    return (false);
  } */
  /*
  if (MsgId == POS_DATA_ID) 
  {
    if (MsgLen < sizeof(T_MsgHdr) + sizeof(T_PositionMsg))
    {
      ErrorMsg  = Origin;
      ErrorMsg += " Client sends insufficient position data, ";
      ErrorMsg += "should be ";
      ErrorMsg += NumToStr (sizeof(T_MsgHdr)+sizeof(T_PositionMsg));
      ErrorMsg += " is: " + NumToStr (MsgHdr->MsgLen);
      AddBadClient (SenderAddress, ErrorMsg, true);
      return false
    }
  }
  */
	return true
} // FG_SERVER::PacketIsValid ()
//////////////////////////////////////////////////////////////////////


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
func (me *FG_SERVER) AddBadClient(Sender *NetAddress, ErrorMsg string, IsLocal bool){
  
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
  me.MaxClientID++
  NewPlayer := new(FG_Player)
  NewPlayer.Callsign      = "* Bad Client *"
  NewPlayer.ModelName     = "* unknown *"
  //NewPlayer.Timestamp     = time(0);
  NewPlayer.JoinTime      = NewPlayer.Timestamp;
 // NewPlayer.Origin        = Sender.Host //getHost ()
  //NewPlayer.Address       = Sender.Address
  NewPlayer.IsLocal       = IsLocal
  NewPlayer.HasErrors     = true
  NewPlayer.Error         = ErrorMsg
  NewPlayer.ClientID      = me.MaxClientID
  NewPlayer.PktsReceivedFrom      = 0
  NewPlayer.PktsSentTo            = 0
  NewPlayer.PktsForwarded         = 0
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


//////////////////////////////////////////////////////////////////////
//  Check if the sender is a known relay, return true if known relay
func (me *FG_SERVER) IsKnownRelay(senderAddress *NetAddress) bool{

  	/*mT_RelayListIt  CurrentRelay = m_RelayList.begin();
  	while (CurrentRelay != m_RelayList.end())
  	{
    	if (CurrentRelay->Address.getIP() == SenderAddress.getIP())
    	{
      	return (true);
    	}
    	CurrentRelay++;
  	}*/
  	_, ok := me.RelayMap[senderAddress.IpAddress]
  	if ok {
  		return true
  	}

  	//string ErrorMsg;
  	//ErrorMsg  = SenderAddress.getHost();
  	//ErrorMsg += " is not a valid relay!";
  	me.AddBlacklist(senderAddress.IpAddress)
  	//SG_LOG (SG_SYSTEMS, SG_ALERT, "UNKNOWN RELAY: " << ErrorMsg);
  	return false
} // FG_SERVER::IsKnownRelay ()




//////////////////////////////////////////////////////////////////////
/**
 * @brief Send any message in m_MessageList to client
 * @param CurrentPlayer Player to send message to
 */
func (me *FG_SERVER) SendChatMessages() {
  
  //mT_MessageIt  CurrentMessage;
	/*
  if ((CurrentPlayer->IsLocal) && (m_MessageList.size()))
  {
    CurrentMessage = m_MessageList.begin();
    while (CurrentMessage != m_MessageList.end())
    {
      if ((CurrentMessage->Target == 0)
      ||  (CurrentMessage->Target == CurrentPlayer->ClientID))
      {
        int len = sizeof(T_MsgHdr) + sizeof(T_ChatMsg);
        m_DataSocket->sendto (CurrentMessage->Msg, len, 0,
          &CurrentPlayer->Address);
      }
      CurrentMessage++;
    }
  } */
} // FG_SERVER::SendChatMessages ()