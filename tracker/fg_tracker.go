

package tracker

const (
	CONNECT = 0
	DISCONNECT =   1
	UPDATE = 2
)

// The Tracker Class
type FG_Tracker struct{
	Server string
	Port int
	Id int // ?
}



// Contruct and return a pointer to a new FG_Tracker
func NewFG_Tracker(server string, port int, id int) *FG_Tracker{
	ob := new(FG_Tracker)
	ob.Server = server
	ob.Port = port
	ob.Id = id// ~what is this
	return ob
}


