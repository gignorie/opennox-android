package opennox

import (
	"os"

	"github.com/opennox/libs/datapath"
	"github.com/opennox/libs/maps"
)

var mapsend struct {
	srv *maps.Server
}

func init() {
	if os.Getenv("NOX_MAPS_HTTP") != "false" {
		registerOnDataPathSet(func() {
			mapsend.srv = maps.NewServer(noxServer.Log, datapath.Maps())
			mapsend.srv.RegisterOnMux(noxServer.HTTP())
		})
	}
}

func clientGetServerMap() string {
	if gameIsNotMultiplayer {
		return nox_xxx_mapFilenameGetSolo_4DB260()
	}
	return noxServer.nox_server_currentMapGetFilename_409B30()
}
