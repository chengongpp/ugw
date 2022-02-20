package ugwrt

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type RtInstance struct {
	Name           string
	WorkDir        string
	Args           []string
	Host           string
	Port           int
	MaxConnections int
	LogLevel       log.Level
	Logger         []*log.Logger
	OutBounds      []OutBound
	Statics        *RtStat
}

type OutBound struct {
	Name     string
	Host     string
	Protocol string
	Port     int
}

type RtStat struct {
	// TODO
}

const (
	AppLog int = iota
	TxLog
	DetailLog
	TraceLog
)

func (rt *RtInstance) Log(logger int, level log.Level, format string, args ...interface{}) {
	if rt.Logger[logger] != nil {
		rt.Logger[logger].Logf(level, format, args...)
	}
}

func (rt *RtInstance) Run() error {
	// Enter the loop
	mainCh := make(chan error, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	defer close(mainCh)
	defer close(sigs)
	go func() {
		// Graceful shutdown
		sig := <-sigs
		fmt.Println(sig)
		mainCh <- nil
	}()
	go rt.MainLoop(mainCh)
	go rt.CtrlLoop(mainCh)

	<-mainCh
	return nil
}

func (rt *RtInstance) MainLoop(mainCh chan error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", rt.Host, rt.Port))
	if err != nil {
		rt.Log(AppLog, log.ErrorLevel, "ListenTCP failed: %v", err)
		mainCh <- err
		return
	}
	rt.Log(AppLog, log.InfoLevel, "ListenTCP on %s:%d", rt.Host, rt.Port)
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			rt.Log(AppLog, log.ErrorLevel, "Close failed: %v", err)
		}
	}(listener)
	// Main loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			rt.Log(AppLog, log.ErrorLevel, "AcceptTCP failed: %v", err)
			mainCh <- err
			continue
		}
		rt.Log(AppLog, log.InfoLevel, "Accepted %s", conn.RemoteAddr().String())
		go rt.HandleConn(conn)
	}
}

func (rt *RtInstance) HandleConn(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			rt.Log(AppLog, log.ErrorLevel, "Close failed: %v", err)
		}
	}(conn)
	// Resolve inbound message
	_, err := rt.ResolveInbound(conn)
	if err != nil {
		rt.Log(AppLog, log.ErrorLevel, "ResolveInbound failed: %v", err)
		return
	}
	// Do custom operations
	// Construct outbound message
	// Send outbound message
}

func (rt *RtInstance) ResolveInbound(conn net.Conn) ([]byte, error) {
	// TODO

	return nil, nil
}

func (rt *RtInstance) CtrlLoop(mainCh chan error) {
	// Accept control commands
	return
}

func (rt *RtInstance) CtrlCmdHandler(signal []byte) error {
	return nil
}
