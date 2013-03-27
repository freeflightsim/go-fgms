

package main

// source = http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/main.cxx
import(
	//"fmt"
	//"strconv"
	"log"
	"io/ioutil"
	"encoding/json"
)
import(
	"github.com/fgx/go-fgms/fgms"
)




// MAIN routine 
func main(){

	// Initialize the beest
	var Servant *fgms.FG_SERVER
	
	Servant = fgms.NewFG_SERVER()
	
	ReadLoadConfigs(Servant, false)

	err := Servant.Init()
	if err != nil {
		//Servant.CloseTracker()
		log.Println("INIT Error", err)
		return
	}
	
	Servant.Loop()
	/*
	I = Servant.Loop();
	if (I != 0)
	{
		Servant.CloseTracker();
		return (I);
	}
	Servant.Done();
	return (0);
	*/
	log.Println("DDDDDDDOOOOOOONEEEE")
}  // main()






// Read a config file and set internal variables accordingly.
func ReadLoadConfigs(Servant *fgms.FG_SERVER, reInit bool) error {

	configFilePath := "/home/gogo/src/github.com/fgx/go-fgms/fgms_example.json"
	
	// Read file
	filebyte, err := ioutil.ReadFile(configFilePath) 
    if err != nil { 
        log.Fatal("Could not read JSON config file: `" + configFilePath + "` ")
        return  err
    } 
    // Parse JSON
	var Config fgms.JSON_ConfAll
    err = json.Unmarshal(filebyte, &Config)
    if err != nil{
    	log.Fatalln("JSON Decode Error from: ", configFilePath,  err)
    	return err
    }
	
	// Server Name
	Servant.SetServerName(Config.Server.Name)
	
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
	
	// Read the list of relays
	for _, relay := range Config.Relays {
		Servant.AddRelay(relay.Host, relay.Port);
	}

	//////////////////////////////////////////////////
	//      read the list of crossfeeds
	//////////////////////////////////////////////////
	Servant.AddCrossfeed ("localhost", 5555)
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
	
	// read the list of blacklisted IPs
	for _, blackList := range Config.Blacklists {
		Servant.AddBlacklist(blackList)
	}
	
	return nil
}

