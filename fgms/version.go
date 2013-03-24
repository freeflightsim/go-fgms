
package fgms

import (
	"fmt"
)

type Version struct{
	Minor int
	Major int
}

func (me *Version) Str() string {
	return fmt.Sprintf("%d.%d", me.Major, me.Minor)
}