

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
	message string
	ch chan string
}


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