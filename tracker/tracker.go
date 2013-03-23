

package tracker


// The Tracker Class
type FG_TRACKER struct{
	Host string
	Port int
	Id int // ?
}

// Contruct and return a pointer to a new FG_TRACKER
func NewFG_TRACKER(host string, port int, id int) *FG_TRACKER{
	ob := new(FG_TRACKER)
	ob.Host = host
	ob.Port = port
	ob.Id = id// ~what is this
	return ob
}