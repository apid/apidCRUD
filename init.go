package apidCRUD

import "github.com/30x/apid-core"

// init() is called by the go runtime.
func init() {
	apid.RegisterPlugin(initPlugin)
}
