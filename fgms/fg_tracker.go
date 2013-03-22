package fgms

// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_tracker.hxx
// http://gitorious.org/fgms/fgms-0-x/blobs/master/src/server/fg_tracker.cxx

// The Tracker Class
type FG_TRACKER struct{
	Host string
	Port int
	Id int
}

// Contruct and return a pointer to a new FG_TRACKER
func NewFG_TRACKER(host string, port int, id int) *FG_TRACKER{
	nu := new(FG_TRACKER)
	nu.Host = host
	nu.Port = port
	nu.Id = id
	return nu
}