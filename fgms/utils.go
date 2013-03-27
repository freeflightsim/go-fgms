
package fgms

import(
	"time"
)

//= Returns an int64 with epoch (should be UTC ?)
func Now() int64{
	return time.Now().Unix()
} 

// TODO There has Got to be a better way
/*
func BytesToString(bites []byte) string{
	cs := ""
	for _, ele := range bites {
		if ele == 0 {
			return cs
		}
		cs += string(ele)
	}
	return cs   
}
*/