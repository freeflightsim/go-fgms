

package fgms

import(
	"net"
	//"bufio"
	//"fmt"
	//"io"
	//"strings"
)

type TcpSrv struct {
	Addr string
}
func (srv *TcpSrv) newConn(rwc net.Conn) (c *TcpConn, err error) {
	c = new(TcpConn)
	c.remoteAddr = rwc.RemoteAddr()
	c.server = srv
	c.rwc = rwc
	return
}

func (srv *TcpSrv) Serve(l net.Listener) error {
	defer l.Close()
	for {
		rw, err := l.Accept()
		if err != nil {
			return err
		}
		c, err := srv.newConn(rw)
		if err != nil {
			continue
		}
		go c.Serve()
	}
	panic("not reached")
}

//==================================
type TcpConn struct {
	remoteAddr	net.Addr
	server		*TcpSrv
	rwc		net.Conn
}
func (c *TcpConn) Serve() {
	var helloString = "hello"
	c.rwc.Write([]byte(helloString))
}

 
// TelnetServer container
type TelnetServer  struct {
	Port int
	Reinit bool
	Received int
	Conn *net.Conn
}

// Constructs and return TelnetServer
func NewTelnetServer() *TelnetServer {
	ob := new(TelnetServer)
	return ob
}


 
// TelnetClient container
type TelnetClient  struct {
	conn net.Conn
	nickname string
	message string
	ch chan string
}

/*
func (c TelnetClient) ReadLinesInto(ch chan TelnetClient) {
	bufc := bufio.NewReader(c.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}
		c.message = fmt.Sprintf("%s", line)
		ch <- c
		
	}
}

func (c TelnetClient) WriteLinesFrom(ch chan TelnetClient) {
	for msg := range ch {
		_, err := io.WriteString(c.conn, c.message)
		if err != nil {
			return
		}
		fmt.Println("Write", msg)
	}
}
*/