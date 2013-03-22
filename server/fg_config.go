package server

// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_config.hxx
// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_config.cxx

import(
	"fmt"
	"io/ioutil"
	"strings"
)

type VarValue struct {
	key string
	val string
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


//////////////////////////////////////////////////////////////////////
/**
 * @brief Read, parse and load a config file 
 * @param ConfigName File name to read
 * @return 1 if an error, else 0
 */
func (me *FG_CONFIG) Read(configFile string) error {

	fmt.Println("READ", configFile )
	//std::ifstream   ConfigFile;
	//std::string     ConfigLine;
	//int             LineNumber;

	//ConfigFile.open (ConfigName.c_str ());
	// Get file contents, return on error
	 contents, err := ioutil.ReadFile(configFile)
	 if err != nil {
	 	return err
	 }
	 //fmt.Println("read", string(contents), err)
	 
	 // Split cotnents into lines
	 lines := strings.Split(string(contents), "\n")
	 fmt.Println("lines=", len(lines))
	 
	//if (!ConfigFile)
	//{
	//	return (1);
	//}
	//LineNumber = 0;
	for idx, line := range lines{
	///while (ConfigFile)
		//fmt.Println(idx, line)
		//getline (ConfigFile, ConfigLine);
		///LineNumber++;
		err := me.ParseLine(line)
		if err != nil {
			fmt.Println("Error in line", idx, line)
		}
		//if (ParseLine (ConfigLine))
		//{
		//	std::cout << "error in line " << LineNumber
		//		<< " in file " << ConfigName
		//		<< std::endl;
		//}
	
	}
	//fmt.Println("REad", me.mT_VarList)
	//ConfigFile.close ();
	//m_CurrentVar = m_VarList.begin ();
	//return (0);
	//*/
	return nil
} // FG_CONFIG::Read ()


//////////////////////////////////////////////////////////////////////
/**
 * @brief Parse the given line, split it into name/value pairs
 *        and put in the internal list
 * @param ConfigLine The line to parse
 * @retval int
 */
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
	me.mT_VarList = append( me.mT_VarList,  &VarValue{key: ki, val: val}  )

	return nil
} // FG_SERVER::ParseLine ()



// Constructs a new isntance
func NewFG_CONFIG() *FG_CONFIG {
	ob := new(FG_CONFIG)
	ob.mT_VarList =  make([]*VarValue,0)
	return ob 
}

