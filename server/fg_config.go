package server

// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_config.hxx
// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_config.cxx

import(
	"fmt"
	"io/ioutil"
	"strings"
)

type VarValue struct {
	Key string
	Val string
}

type FG_CONFIG struct {

	//public:
	//int   Read ( const std::string &ConfigName );
	//void  Dump ();
	//void  SetStart ();
	//int   Next ();
	//std::string Get ( const std::string &VarName );
	//std::string GetName ();
	//std::string GetValue ();
	//std::string GetNext ();
	//int SetSection ( const std::string &SecName );
	//int SecNext ();
	//std::string GetSecNextVar ();
	//std::string GetSecNextVal ();
	//private:
	//typedef std::pair<std::string,std::string>  mT_VarValue;
	//typedef std::list<mT_VarValue>              mT_VarList;
	//int   ParseLine ( const std::string &ConfigLine );
	//mT_VarList            m_VarList;
	mT_VarList []*VarValue
	//mT_VarList::iterator  m_CurrentVar;
	//std::string           m_CurrentSection;
}

// Constructs a new isntance
func NewFG_CONFIG() *FG_CONFIG {
	ob := new(FG_CONFIG)
	ob.mT_VarList =  make([]*VarValue,0)
	return ob 
}

//////////////////////////////////////////////////////////////////////
/**
 * @brief Read, parse and load a config file 
 * @param ConfigName File name to read
 * @return 1 if an error, else 0
 */
func (me *FG_CONFIG) Read(configFile string) error {

	fmt.Println("READ", configFile )
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
} // FG_CONFIG::Read ()


//////////////////////////////////////////////////////////////////////
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

	// Find position of first =
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
} // FG_SERVER::ParseLine ()






// Find a variable with name 'VarName' in the internal list and return
 
func (me *FG_CONFIG) Get(VarName string) string {
	for _, ele := range me.mT_VarList {
		if ele.Key == VarName {
			return ele.Val
		}
	}
	return ""
	/* m_CurrentVar = m_VarList.begin();
	while (m_CurrentVar != m_VarList.end())
	{
		if (m_CurrentVar->first == VarName)
		{
			return (m_CurrentVar->second);
		}
		m_CurrentVar++;
	}
	return (""); */
} // FG_SERVER::Get ()
