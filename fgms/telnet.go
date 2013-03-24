

package fgms

import(
	"net"
	"bufio"
	"fmt"
	"io"
	//"strings"
)
 
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
	ch chan string
}


func (c TelnetClient) ReadLinesInto(ch chan<- string) {
	bufc := bufio.NewReader(c.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}
		ch <- fmt.Sprintf("%s: %s", c.nickname, line)
	}
}

func (c TelnetClient) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		_, err := io.WriteString(c.conn, msg)
		if err != nil {
			return
		}
	}
}