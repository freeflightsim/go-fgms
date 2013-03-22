

package main

// source = http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/main.cxx
import(
	"fmt"
)
import(
	"github.com/fgx/go-fgms/server"
)


//////////////////////////////////////////////////////////////////////
/// MAIN routine 
//* @param argc
// @param argv*[]
//int main ( int argc, char* argv[] )

func main(){

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
//////////////////////////////////////////////////////////////////////



/////////////////////////////////////////////////////////////////////
/// (re)Read config files - ReInit True to reinitialize
//int ReadConfigs ( bool ReInit = false )

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


/**
 * @brief Read a config file and set internal variables accordingly
 * @param ConfigName Path of config file to load
 * @retval int  -- todo--
 */
//bool ProcessConfig ( const string& ConfigName )
func ProcessConfig( configFilePath string) error{

	Config := new(server.FG_CONFIG)
	//var Val string
	//var E int

	//if (bHadConfig)	// we already have a config, so ignore
	//	return (true);
	err := Config.Read (configFilePath)
	//if (ok){
	//fmt.Println("err=", err);
	if err != nil {
		return err
	}
	//SG_ALERT (SG_SYSTEMS, SG_ALERT, "processing " << ConfigName);
	//Val = Config.Get ("server.name");
	//if (Val != "")
	//{
	//	Servant.SetServerName (Val);
	//	bHadConfig = true; // got a serve name - minimum 
	//}
	/*Val = Config.Get ("server.address");
	if (Val != "")
	{
		Servant.SetBindAddress (Val);
	}
	Val = Config.Get ("server.port");
	if (Val != "")
	{
		Servant.SetDataPort (StrToNum<int> (Val.c_str (), E));
		if (E)
		{
			SG_ALERT (SG_SYSTEMS, SG_ALERT, "invalid value for DataPort: '" << optarg << "'");
			exit (1);
		}
	}
	Val = Config.Get ("server.telnet_port");
	if (Val != "")
	{
		Servant.SetTelnetPort (StrToNum<int> (Val.c_str (), E));
		if (E)
		{
			SG_ALERT (SG_SYSTEMS, SG_ALERT, "invalid value for TelnetPort: '" << optarg << "'");
			exit (1);
		}
	}
	Val = Config.Get("server.out_of_reach");
	if (Val != "")
	{
		Servant.SetOutOfReach (StrToNum<int> (Val.c_str (), E));
		if (E)
		{
			SG_ALERT (SG_SYSTEMS, SG_ALERT, "invalid value for OutOfReach: '" << optarg << "'");
			exit (1);
		}
	}
	Val = Config.Get("server.playerexpires");
	if (Val != "")
	{
		Servant.SetPlayerExpires (StrToNum<int> (Val.c_str (), E));
		if (E)
		{
			SG_ALERT (SG_SYSTEMS, SG_ALERT, "invalid value for Expire: '" << optarg << "'");
			exit (1);
		}
	}
	Val = Config.Get ("server.logfile");
	if (Val != "")
	{
		Servant.SetLogfile (Val);
	}
	Val = Config.Get ("server.daemon");
	if (Val != "")
	{
		if ((Val == "on") || (Val == "true"))
		{
			RunAsDaemon = true;
		}
		else if ((Val == "off") || (Val == "false"))
		{
			RunAsDaemon = false;
		}
		else
		{
			SG_ALERT (SG_SYSTEMS, SG_ALERT, "unknown value for 'server.daemon'!" << " in file " << ConfigName);
		}
	}
	Val = Config.Get ("server.tracked");
	if (Val != "")
	{
		string  Server;
		int     Port;
		bool    tracked;
		if (Val == "true")
		{
			tracked = true;
		}
		else
		{
			tracked = false;
		}
		Server = Config.Get ("server.tracking_server");
		Val = Config.Get ("server.tracking_port");
		Port = StrToNum<int> (Val.c_str (), E);
		if (E)
		{
			SG_ALERT (SG_SYSTEMS, SG_ALERT, "invalid value for tracking_port: '" << Val << "'");
			exit (1);
		}
    if ( tracked && ( Servant.AddTracker (Server, Port, tracked) != FG_SERVER::SUCCESS ) ) // set master m_IsTracked
    {
			SG_ALERT (SG_SYSTEMS, SG_ALERT, "Failed to get IPC msg queue ID! error " << errno );
			exit (1); // do NOT continue if a requested 'tracker' FAILED
    }
	}
	Val = Config.Get ("server.is_hub");
	if (Val != "")
	{
		if (Val == "true")
		{
			Servant.SetHub (true);
		}
		else
		{
			Servant.SetHub (false);
		}
	}
	//////////////////////////////////////////////////
	//      read the list of relays
	//////////////////////////////////////////////////
	bool    MoreToRead  = true;
	string  Section = "relay";
	string  Var;
	string  Server = "";
	int     Port   = 0;
	if (! Config.SetSection (Section))
	{
		MoreToRead = false;
	}
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
	//////////////////////////////////////////////////
	//      read the list of crossfeeds
	//////////////////////////////////////////////////
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
	//////////////////////////////////////////////////
	//      read the list of blacklisted IPs
	//////////////////////////////////////////////////
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

