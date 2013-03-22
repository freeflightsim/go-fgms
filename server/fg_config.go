package server

// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_config.hxx
// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_config.cxx

import(
	"fmt"
	"io/ioutil"
	"strings"
)

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
	 fmt.Println("read", string(contents), err)
	 
	 // Split cotnents into lines
	 lines := strings.Split(string(contents), "\n")
	 fmt.Println("lines=", len(lines))
	 
	//if (!ConfigFile)
	//{
	//	return (1);
	//}
	//LineNumber = 0;
	/*while (ConfigFile)
	{
		getline (ConfigFile, ConfigLine);
		LineNumber++;
		if (ParseLine (ConfigLine))
		{
			std::cout << "error in line " << LineNumber
				<< " in file " << ConfigName
				<< std::endl;
		}
	}
	ConfigFile.close ();
	m_CurrentVar = m_VarList.begin ();
	return (0);
	*/
	return nil
} // FG_CONFIG::Read ()

