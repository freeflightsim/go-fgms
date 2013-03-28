

package fgms

import(
	//"encoding/json"
	//"ioutil"
	//"log"
	
)
//= A Host Row
type JSON_HostConf struct{
	Host string `json:"host"`
	Port int `json:"port"`
	Comment string `json:"comment"`
}

//= Main Server Configuration
type JSON_ServerConf struct{
	
	Name string `json:"name"`
	Address string `json:"address"`
	Port int `json:"port"`
	IsHub bool `json:"is_hub"`
	
	LogFile string `json:"log_file"`
	
	TelnetPort int `json:"telnet_port"`
	
	PlayerExpiresSecs int `json:"player_expires"`
	OutOfReachNm int `json:"out_of_reach"`
	
	Tracked bool `json:"tracked"`
	TrackingServer string `json:"tracking_server"`
		
}
//= Whole Payload from File
type JSON_ConfAll struct {
	Server JSON_ServerConf `json:"server"`
	Relays []JSON_HostConf `json:"relays"`
	Crossfeeds []JSON_HostConf	 `json:"crossfeeds"`
	Blacklists []string		 `json:"blacklists"`
}


