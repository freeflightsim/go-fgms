

package main

// source = http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/main.cxx
import(
	"fmt"
	"strconv"
)
import(
	"github.com/fgx/go-fgms/fgms"

)

// Main instance of FG_SERVER
var Servant *fgms.FG_SERVER

//////////////////////////////////////////////////////////////////////
/// MAIN routine 
//* @param argc
// @param argv*[]
//int main ( int argc, char* argv[] )
func main(){

	// Initialize the beest
	Servant = fgms.NewFG_SERVER()
	//int     I;
	//#if defined ENABLE_DEBUG
	//  logbuf::set_log_classes(SG_GENERAL);
	//#endif
	
	//ParseParams (argc, argv);
	ReadConfigs(false)
	//if ( !bHadConfig )
	//{
	//	SG_ALERT (SG_SYSTEMS, SG_ALERT, "No configuration file '" << DEF_CONF_FILE << "' found!");
	//	exit(1);
	//}
	/*
	sglog().setLogLevels( SG_ALL, SG_INFO );
	sglog().enable_with_date (true);
	I = Servant.Init ();
	if (I != 0)
	{
		Servant.CloseTracker();
		return (I);
	}
	#ifndef _MSC_VER
	if (RunAsDaemon)
	{
		Myself.Daemonize ();
		SG_ALERT (SG_SYSTEMS, SG_ALERT, "Main server started!");
	}
	#endif
	I = Servant.Loop();
	if (I != 0)
	{
		Servant.CloseTracker();
		return (I);
	}
	Servant.Done();
	return (0);
	*/
}  // main()




// (re)Read config files - ReInit True to reinitialize
func ReadConfigs(ReInit bool) error {
	
	var Path string
	Path = "/home/gogo/src/github.com/fgx/go-fgms/fgms_example.conf"
	
	//if (Path != "")
	//{
	//	Path += "/" DEF_CONF_FILE;
	err := ProcessConfig (Path)
	if err != nil {
		fmt.Println("error=", err)
		return err
	}
	//if (ProcessConfig (DEF_CONF_FILE))
	//	return 1;
	return nil
} // ReadConfigs ()


// Read a config file and set internal variables accordingly.
// Returns an error or nil
func ProcessConfig( configFilePath string) error{

	Config := fgms.NewFG_CONFIG()
	
	//if (bHadConfig)	// we already have a config, so ignore
	//	return (true);
	err := Config.Read (configFilePath)
	//if (ok){
	//fmt.Println("err=", err);
	if err != nil {
		return err
	}
	//SG_ALERT (SG_SYSTEMS, SG_ALERT, "processing " << ConfigName);
	//var Val server.VarValue
	
	// Server Name
	Val := Config.Get("server.name")
	fmt.Println("server.name", Val)
	if Val != "" {
		Servant.SetServerName(Val)
		//	bHadConfig = true; // got a serve name - minimum 
	}
	
	// Address
	Val = Config.Get("server.address")
	fmt.Println("server.address", Val)
	if Val != "" {
		Servant.SetBindAddress(Val)
	}
	
	// UDP Port No
	Val = Config.Get("server.port")
	fmt.Println("server.port", Val)
	if Val != "" {
		port, err := strconv.ParseInt(Val, 10, 0)
		if err != nil {
			fmt.Println("Error", "invalid value for DataPort")
			return err
		} 
		Servant.SetDataPort(int(port))
	}
	
	// Telnet Port
	Val = Config.Get ("server.telnet_port");
	if Val != "" {
		telnetport, err := strconv.ParseInt(Val, 10, 0)
		if err != nil {
			fmt.Println("Error", "invalid value for Telnet")
			return err
		} 
		Servant.SetTelnetPort(int(telnetport))
	}
	
	// Outta Reach
	nm, err := Config.GetInt("server.out_of_reach")
	if err != nil {
		fmt.Println("Error", "invalid value for `server.out_of_reach`", Val)
		return err	
	}
	Servant.SetOutOfReach( nm )

	// Player Expires	
	exp_secs, err := Config.GetInt("server.playerexpires")
	if err != nil {
		fmt.Println("Error", "invalid value for `server.playerexpires`", err, Val)
		return err	
	}
	Servant.SetPlayerExpires( exp_secs )
	
	// Server is hub
	Val = Config.Get("server.is_hub")
	if Val != "" {
		is_hub, err := strconv.ParseBool(Val)
		if err != nil {
			fmt.Println("Error", "server.is_hub", Val)
			return err
		}
		Servant.SetHub( is_hub ) 
	}
	
	
	
	// Log File
	Val = Config.Get("server.logfile")
	if Val != "" {
		Servant.SetLogfile(Val);
	}
	
	// Tracked
	Val = Config.Get ("server.tracked")
	if Val != "" {
		tracked, _ := strconv.ParseBool(Val)
		if tracked {
			trkServer := Config.Get("server.tracking_server")
			trkPorts := Config.Get("server.tracking_port")
			trkPorti, err := strconv.ParseInt(trkPorts, 10, 0)
			if err != nil{
				fmt.Println("Error", "invalid value for tracking_port: ", Val)
				return err
			}
			pii := int(trkPorti)
			fmt.Println("addd", trkServer, pii, tracked)
			Servant.AddTracker(trkServer, pii, tracked)
			
		} 
	}
	
	

	
	//////////////////////////////////////////////////
	//      read the list of relays
	//////////////////////////////////////////////////
	
	//bool    MoreToRead  = true;
	//Section := "relay"
	//string  Var;
	//string  Server = "";
	//int     Port   = 0;
	//if (! Config.SetSection (Section))
	//{
	//	MoreToRead = false;
	//}
	//= not sure how this works in relay
	vals, err := Config.GetSection("relay")
	if err != nil {
		fmt.Println("section not found")
		return err
	}
	
	if len(vals) > 0 {
		fmt.Println("section=", vals)
		server := vals["relay.host"]
		port, err := Config.GetInt("relay.port") 
		if err != nil{
			fmt.Println("Error:", err)
			return err
		}
		Servant.AddRelay(server, port);
	}
	/*
	while (MoreToRead)
	{
		Var = Config.GetName ();
		Val = Config.GetValue();
		if (Var == "relay.host")
		{ 
			Server = Val;
		}
		if (Var == "relay.port")
		{ 
			Port = StrToNum<int> (Val.c_str(), E);
			if (E)
			{ 
				SG_ALERT (SG_SYSTEMS, SG_ALERT, "invalid value for RelayPort: '" << Val << "'");
				exit (1);
			}
		}
		if ((Server != "") && (Port != 0))
		{ 
			Servant.AddRelay (Server, Port);
			Server = "";
			Port   = 0;
		}
		if (Config.SecNext () == 0)
		{ 
			MoreToRead = false;
		}
	}
	*/
	//////////////////////////////////////////////////
	//      read the list of crossfeeds
	//////////////////////////////////////////////////
	/*
	MoreToRead  = true;
	Section = "crossfeed";
	Var    = "";
	Server = "";
	Port   = 0;
	if (! Config.SetSection (Section))
	{
		MoreToRead = false;
	}
	while (MoreToRead)
	{
		Var = Config.GetName ();
		Val = Config.GetValue();
		if (Var == "crossfeed.host")
		{
			Server = Val;
		}
		if (Var == "crossfeed.port")
		{
			Port = StrToNum<int> (Val.c_str(), E);
			if (E)
			{
				SG_ALERT (SG_SYSTEMS, SG_ALERT, "invalid value for crossfeed.port: '" << Val << "'");
				exit (1);
			}
		}
		if ((Server != "") && (Port != 0))
		{
			Servant.AddCrossfeed (Server, Port);
			Server = "";
			Port   = 0;
		}
		if (Config.SecNext () == 0)
		{
			MoreToRead = false;
		}
	}
	*/
	//////////////////////////////////////////////////
	//      read the list of blacklisted IPs
	//////////////////////////////////////////////////
	blacklist, err := Config.GetList("blacklist")
	fmt.Println(blacklist, err)
	if len(blacklist) > 0 {
		for _, bl := range blacklist {
			Servant.AddBlacklist(bl)
		}
	}
	/*
	MoreToRead  = true;
	Section = "blacklist";
	Var    = "";
	Val    = "";
	if (! Config.SetSection (Section))
	{
		MoreToRead = false;
	}
	while (MoreToRead)
	{
		Var = Config.GetName ();
		Val = Config.GetValue();
		if (Var == "blacklist")
		{
			Servant.AddBlacklist (Val);
		}
		if (Config.SecNext () == 0)
		{
			MoreToRead = false;
		}
	}*/
	//////////////////////////////////////////////////
	return nil
} // ProcessConfig ( const string& ConfigName )

