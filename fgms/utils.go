
package fgms

import(
	"fmt"
	"time"
)
 import(
	"github.com/FreeFlightSim/go-fgms/flightgear"
)
//= Returns an int64 with epoch (should be UTC ?)
func Now() int64{
	//time.UTC().unix() ??
	return time.Now().Unix()
} 

// pete FAIL FAIL FAIL's 
func GetProtocolVersionString() string {
	//return "1.1"
	major := flightgear.PROTO_VER >> 16
	minor := flightgear.PROTO_VER & 0xffff
	return fmt.Sprintf("%d.%d", major, minor) 
}


// TODO There has Got to be a better way
/*
func BytesToString(bites []byte) string{
	cs := ""
	for _, ele := range bites {
		if ele == 0 {
			return cs
		}
		cs += string(ele)
	}
	return cs   
}
*/