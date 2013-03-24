
package fgms

import (
	"fmt"
)

// Version Flags
type Version struct{
	Minor int
	Major int
}

// Nice string "n.n"
func (me *Version) Str() string {
	return fmt.Sprintf("%d.%d", me.Major, me.Minor)
}