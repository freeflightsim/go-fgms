
package message

import (
	"errors"
	//"fmt"
	//"time"

	"github.com/davecgh/go-xdr/xdr"
)


var ErrDecode = errors.New("message: XDR decode error")
var ErrProtoVer = errors.New("message: Invalid protocol version")
var ErrMagic = errors.New("message: Invalid magic")

// Decode the Header part of the xdr_encoded bytes,
// returns the header, remainder bytes and error
func DecodeHeader(xdr_enc []byte)(HeaderMsg, []byte, error) {


	var header HeaderMsg

	remainingBytes, err := xdr.Unmarshal(xdr_enc, &header)
	if err != nil{
		return header, remainingBytes, ErrDecode
	}

	if header.Version != PROTOCOL_VER {
		return header, remainingBytes, ErrProtoVer
	}
	/*
	switch header.Magic {
		case MSG_MAGIC:  fallthrough
		case RELAY_MAGIC:  fallthrough
		default:
			return header, remainingBytes, errors.New("Invalid Magic")
	}
	*/
	if header.Magic != MSG_MAGIC && header.Magic != RELAY_MAGIC {
		return header, remainingBytes, ErrMagic
	}
	//if header.Type != TYPE_POS {
	//	return header, errors.New("Not a position error")
	//}

	return header, remainingBytes, nil
}
