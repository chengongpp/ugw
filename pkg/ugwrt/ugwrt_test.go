package ugwrt

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestUgwRuntime(t *testing.T) {
	//TODO
	rt := &RtInstance{
		Name:           "test",
		WorkDir:        "",
		Args:           os.Args,
		Host:           "",
		Port:           0,
		MaxConnections: 0,
		LogLevel:       log.DebugLevel,
		Logger:         nil,
		OutBounds:      nil,
		Statics:        nil,
	}
}
