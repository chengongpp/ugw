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
		Outbounds:      nil,
		Stats:          nil,
	}
	rt.Run()
}

func TestResolveProtos(t *testing.T) {
	config := Config{
		LogLevel: "debug",
	}
	rt := NewRtInstance(config)
	_, err := rt.ResolveProtos("C:\\Users\\cheng\\IdeaProjects\\ugw\\layout/messages/")
	if err != nil {
		t.Error("ResolveProtos failed")
	}
}
