
package message

import (
	"errors"
	"fmt"
	//"time"

	"github.com/davecgh/go-xdr/xdr"
)




// Decode the Header part of packet
func DecodeHeader(xdr_enc []byte)(HeaderMsg, []byte, error) {

	var header HeaderMsg

	remainingBytes, err := xdr.Unmarshal(xdr_enc, &header)
	if err != nil{
		fmt.Println("XDR Decode Error", err)
		return header, remainingBytes, nil
	}
	//fmt.Println("remain=", len(remainingBytes))
	//fmt.Println( header.Magic == MSG_MAGIC, header.Version ==  PROTOCOL_VER)
	//fmt.Println ("Header=", len(remainingBytes), header.Type, header.Type == TYPE_POS, header.Version, header.Callsign(), )

	if header.Version != PROTOCOL_VER {
		return header, remainingBytes, errors.New("Invalid protocol version")
	}
	switch header.Magic {
		case MSG_MAGIC:  fallthrough
		case RELAY_MAGIC:  fallthrough
		default:
			return header, remainingBytes, errors.New("Invalid Magic")
	}

	//if header.Type != TYPE_POS {
	//	return header, errors.New("Not a position error")
	//}

	/*
	var position PositionMsg
	rembits, err := xdr.Unmarshal(remainingBytes, &position)

	if err != nil {
		fmt.Println(rembits)
	}
	t := time.Unix(0, int64(position.Time) * int64(time.Nanosecond) ).UTC()
	t2 := time.Unix(0, int64(position.Time) * int64(time.Millisecond))
	if 1 == 2 {
		fmt.Println(position.Lag, position.Time, ">>", t, "==", t2, header.Callsign(), ":", len(rembits))
	}
	*/
	return header, remainingBytes, nil
}
