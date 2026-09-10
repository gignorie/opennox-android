package opennox

import (
	"github.com/opennox/opennox/v1/legacy/timer"
)

func init() {
	timer.PlatformTicks = platformTicks
}
