
package fgms

import(
	"fmt"
	"log"
	"net"
)

import(
	"github.com/davecgh/go-xdr/xdr"
	"github.com/fgx/go-fgms/flightgear"
)
//  Insert a new crossfeed server into internal list - after resolution of address
func (me *FG_SERVER) AddCrossfeed( host_name string, port int){

	log.Println("> Add Crossfeed = ", host_name, port)

	//= Now go and do check is background
	go func(host_name string, port int){
		
		//= Now resolve with UDP address			
		host_port := fmt.Sprintf("%s:%d", host_name, port)
		udp_addr, err := net.ResolveUDPAddr("udp", host_port)
		if err != nil {
			log.Println("\tFAIL: Crossfeed to resolve UDP address:  ", host_port, err)
			return
		}
		
		//= Now we open socket and listen
		var err_listen error
		me.Crossfeeds[host_port], err_listen = net.ListenUDP("udp", udp_addr)
		if err_listen != nil {
			log.Println("\tFAIL: Crossfeed FAIL to Open:  ", host_port, udp_addr, err_listen)
			return
		}
		log.Println("\tOK:   Crossfeed Added -  ", host_port, udp_addr, err_listen)
		
	}(host_name, port)	
} // FG_SERVER::AddCrossfeed()




/*   Send message to all crossfeed servers.
	Crossfeed servers receive all traffic without condition,
	mainly used for testing and debugging, and crossfeed.fgx.ch
	http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_server.cxx#line1154
*/
func (me *FG_SERVER) SendToCrossfeed(Msg []byte, Bytes int, SenderAddress *net.UDPAddr){

	//T_MsgHdr*       MsgHdr;
	//uint32_t        MsgMagic;
	//int             sent;
	
	//MsgHdr    = (T_MsgHdr *) Msg;
	//MsgMagic  = MsgHdr->Magic;
	//MsgHdr->Magic = XDR_encode<uint32_t> (RELAY_MAGIC);
	
	// Not sure what is happening, but we create another payload
	var MsgHdr flightgear.T_MsgHdr
	_, err := xdr.Unmarshal(Msg, &MsgHdr)
	if err != nil {
		fmt.Println("XDR Decode Error in SendToCrossfeed - Should never happen?", err)
		return
	}
	MsgHdr.Magic = RELAY_MAGIC
	
	encoded, err := xdr.Marshal(MsgHdr)
	if err != nil {
		fmt.Println("XDR Encode Error in SendToCrossfeed - Should never happen?", err)
		return
	}
	//mT_RelayListIt CurrentCrossfeed = m_CrossfeedList.begin();
	//while (CurrentCrossfeed != m_CrossfeedList.end())
	//{
	for _, loopCF := range me.Crossfeeds {

		//if (CurrentCrossfeed->Address.getIP() != SenderAddress.getIP()) ?? But same address different port ??
		_, err := loopCF.WriteToUDP(encoded, SenderAddress)
		if err != nil {
			me.CrossFeedFailed++
		}else {
			me.CrossFeedSent++
		}
		//CurrentCrossfeed++;
	}
	//MsgHdr->Magic = MsgMagic;  // restore the magic value ? umm not used now ?
} // FG_SERVER::SendToCrossfeed ()

