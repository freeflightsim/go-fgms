package fgms

// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_config.hxx
// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_config.cxx

import(
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"errors"
)

// A key/val pair
type VarValue struct {
	Key string
	Val string
}

// The Config
type FG_CONFIG struct {
	mT_VarList []*VarValue
}

// Construct and return a new instance if FG_CONFIG
func NewFG_CONFIG() *FG_CONFIG {
	ob := new(FG_CONFIG)
	ob.mT_VarList =  make([]*VarValue,0)
	return ob 
}



// Read, parse and load a config file
func (me *FG_CONFIG) Read(configFile string) error {

	// Get file contents, return on error
	contents, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	
	// Split cotnents into lines
	lines := strings.Split(string(contents), "\n")
	
	// loops the lines and parse into list 
	for idx, line := range lines{
		err := me.ParseLine(line)
		if err != nil {
			fmt.Println("Error in line", idx, line)
		}
	}
	return nil
} 


// Parse the given line, split it into name/value pairs
// and put in the internal list
func(me *FG_CONFIG) ParseLine(ConfigLine string) error{

	// Strips whitespace 
	ConfigLine = strings.TrimSpace(ConfigLine)
	
	if len(ConfigLine) == 0 {
		return nil // Blank line
	}
	if ConfigLine[0] == '#' {
		return nil // # comments line
	}

	// Find position of first "="
	eqPos := strings.IndexAny( ConfigLine, "=")
	if eqPos < 2 {
		return nil // less than is silly atmo
	}
	// Get the ki and val and stip them of whitespace	
	ki := strings.TrimSpace( ConfigLine[: eqPos -1 ] )
	val := strings.TrimSpace( ConfigLine[eqPos +1 :] )
	
	// Add to out list
	me.mT_VarList = append( me.mT_VarList,  &VarValue{Key: ki, Val: val}  )

	return nil
} 



// Find a variable with name 'VarName' in the internal list and return its string value
func (me *FG_CONFIG) Get(VarName string) string {
	for _, ele := range me.mT_VarList {
		if ele.Key == VarName {
			return ele.Val
		}
	}
	return ""
}

// Find a variable with name 'VarName' in the internal list and return its int value
func (me *FG_CONFIG) GetInt(VarName string) (int, error) {
	for _, ele := range me.mT_VarList {
		if ele.Key == VarName {
			i, err := strconv.ParseInt(ele.Val, 10, 0)
			return int(i), err
		}
	}
	return 0, errors.New("Cant find Var")
}

// Find a variable with name 'VarName' in the internal list and return its int value
func (me *FG_CONFIG) GetList(VarName string) ([]string, error) {
	lst := make([]string, 0)
	for _, ele := range me.mT_VarList {
		if ele.Key == VarName {
			lst = append(lst, ele.Val)
			
		}
	}
	return lst, nil
}

// Find a variable with name 'VarName' in the internal list and return bool value
/*
 func (me *FG_CONFIG) GetBool(VarName string) (bool error) {
	for _, ele := range me.mT_VarList {
		if ele.Key == VarName {
			return strconv.ParseBool(ele.Val)			
		}
	}
	return nil, nil
}
*/


// Return a list of the matching to section. 
/* 
   Example input line
		relay.host = mpserver01.flightgear.org
		relay.port = 5000
	GetSection("relay")
	returns map string/string
	   relay.host =  mpserver01.flightgear.org
	   relay.port = 500
*/
 func (me *FG_CONFIG) GetSection(sec_name string) (map[string]string, error) {

	// Section name stats with foo eg foo.
 	sec_name += "." 
 	lenny := len(sec_name) // length which we use to slice
 	ret := make(map[string]string)
	for _, ele := range me.mT_VarList {
		ki := ele.Key[:lenny] // try and find foo.
		if ki == sec_name { 
			val := me.Get(ele.Key)
			ret[ele.Key] = val
			fmt.Println(">", sec_name, lenny, ki, val)
		}
	}
	return ret, nil
}

