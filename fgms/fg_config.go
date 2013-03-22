package fgms

// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_config.hxx
// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_config.cxx

import(
	"fmt"
	"io/ioutil"
	"strings"
	//"strconv"
)

type VarValue struct {
	Key string
	Val string
}

type FG_CONFIG struct {
	mT_VarList []*VarValue
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



// Find a variable with name 'VarName' in the internal list and return its value
 func (me *FG_CONFIG) Get(VarName string) string {
	for _, ele := range me.mT_VarList {
		if ele.Key == VarName {
			return ele.Val
		}
	}
	return ""
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

// Construct and return a new instance if FG_CONFIG
func NewFG_CONFIG() *FG_CONFIG {
	ob := new(FG_CONFIG)
	ob.mT_VarList =  make([]*VarValue,0)
	return ob 
}
