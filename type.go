// webserver
package webserver

import (
	"github.com/fe0b6/webobj"
)

type InitObj struct {
	Port  int
	Route func(*webobj.RqObj)
}
