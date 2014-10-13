

package main

// source = http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/main.cxx
import(
	//"fmt"
	//"strconv"
	"flag"
	"log"
	//"io/ioutil"
	//"encoding/json"

	"github.com/freeflightsim/go-fgms/fgms"
)




// MAIN routine 
func main(){

	var iconfig *string = flag.String("c", "./fgms_example.json", "Path to config file")
	flag.Parse()
	

	
	//ReadLoadConfigs(Servant, false)
	config, err := fgms.LoadConfig(*iconfig)
	if err != nil {
		log.Println("Cannot load config")
		return
	}

	// Initialize the beest
	//var server *fgms.FgServer

	//server = fgms.NewFgServer()
	log.Println("main=", fgms.Server)
	fgms.SetConfig(config)

	//err_init := server.Init()
	//if err_init != nil {
		//Servant.CloseTracker()
	//	log.Println("INIT Error", err_init)
	//	return
	//}

	fgms.Server.Start()
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






