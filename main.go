

package main

// source = http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/main.cxx
import(
	"fmt"
	//"strconv"
	"log"
	"io/ioutil"
	"encoding/json"
)
import(

	//"github.com/kylelemons/go-gypsy/yaml"

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
	
	//sglog().setLogLevels( SG_ALL, SG_INFO );
	//sglog().enable_with_date (true);
	err := Servant.Init()
	if err != nil {
		//Servant.CloseTracker()
		return
	}
	/*
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
	
	//yamlFile :=  "/home/gogo/src/github.com/fgx/go-fgms/fgms_example.yaml"
	//conf, erry := yaml.ReadFile(yamlFile)
	
	//fmt.Println("error=", conf, erry)
	
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

	//yamlFile :=  "/home/gogo/src/github.com/fgx/go-fgms/fgms_example.yaml"
	jsonFile :=  "/home/gogo/src/github.com/fgx/go-fgms/fgms_example.json"
	
	//Config, errj := LoadJsonConfig(jsonFile)
	filebyte, err := ioutil.ReadFile(jsonFile) 
    if err != nil { 
        log.Fatal("Could not read file " + jsonFile + " to parse")
        return  err
    } 
    //log.Println(string(filebyte))
	var Config fgms.JSON_ConfAll
    err = json.Unmarshal(filebyte, &Config)
    if err != nil{
    	fmt.Println("JSON Decode Error", err)
    	return err
    }
	
	Servant.SetServerName(Config.Server.Name)
	//	bHadConfig = true; // got a serve name - minimum
	//return nil
	
	// Address
	Servant.SetBindAddress(Config.Server.Address)
	
	// UDP Port No
	Servant.SetDataPort(Config.Server.Port)
	
	
	// Telnet Port
	Servant.SetTelnetPort(Config.Server.TelnetPort)
	
	// Outta Reach
	Servant.SetOutOfReach(Config.Server.OutOfReachNm)
	
	// Player Expires
	Servant.SetPlayerExpires(Config.Server.PlayerExpiresSecs)	
	
	// Server is hub
	Servant.SetHub( Config.Server.IsHub ) 

	
	// Log File
	Servant.SetLogfile(Config.Server.LogFile);
	
	// Tracked
	/*
	Val, err = Config.Get ("server.tracked")
	if Val != "" {
		tracked, _ := strconv.ParseBool(Val)
		if tracked {
			trkServer, err := Config.Get("server.tracking_server")
			if err != nil {
				log.Fatalln("Error", "Missing `server.tracking_server`", trkServer)
				return err
			}
			fmt.Println("TRK", trkServer,  tracked)
			Servant.AddTracker(trkServer, pii, tracked)
			
		} 
	}
	*/
	

	
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
	//relays, err := Config.Get("relays")
	//log.Fatalln("Error", "No  `relays` found", relays, err)
	//if err != nil {
		//fmt.Println("section not found")
	//	log.Fatalln("Error", "No  `relays` found", Val)
	//	return err
	//}
	//fmt.Println("RELAYS:", relays)
	for idx, ele := range Config.Relays {
		fmt.Println("REPLAY=", idx, ele)
	
	}
	/* if len(vals) > 0 {		
		server := vals["relay.host"]
		port, err := Config.GetInt("relay.port") 
		if err != nil{
			//fmt.Println("Error:", err)
			return err
		}
		Servant.AddRelay(server, port);
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
	/* blacklist, err := Config.GetList("blacklist")
	fmt.Println(blacklist, err)
	if len(blacklist) > 0 {
		for _, bl := range blacklist {
			Servant.AddBlacklist(bl)
		}
	} */
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

