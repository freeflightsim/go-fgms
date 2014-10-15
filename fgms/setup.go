
package fgms


// Auto Initialize App
func init(){


	SetupBlackList()
	SetupCrossfeed()
	SetupRelays()
	SetupHttp()

	// Server must be last
	SetupServer()
}
