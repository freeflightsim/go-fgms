
package fgms

//= A struct for general statistics usage and incrementals
type Stats struct{
	Recv 	uint64
	Failed 	uint64
	Sent	uint64
	Invalid uint64
	Forward uint64
}
