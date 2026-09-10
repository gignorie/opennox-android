package version

import (
	"github.com/opennox/libs/log"
)

var Log = log.New("version")

func LogVersion() {
	Log.Info("engine", "version", Version(), "commit", Commit())
}

const (
	DefVersion = "v1.9.0-dev"
)

const (
	devCommit = "<dev>"
)

var (
	version = DefVersion
	commit  = devCommit
)

func Version() string {
	return version
}

func Commit() string {
	return commit
}

func ClientVersion() string {
	if IsDev() {
		return version + " (" + commit + ")"
	}
	return version
}

func IsDev() bool {
	return commit == devCommit || semverIsDev(version)
}
