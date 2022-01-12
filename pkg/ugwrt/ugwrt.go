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
		rt.Logger[logger].Printf(format, args...)
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
	defer listener.Close()
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

func (rt *RtInstance) HandleConn(conn net.Conn) error {
	return nil
}

func (rt *RtInstance) CtrlLoop(mainCh chan error) error {
	return nil
}

func (rt *RtInstance) CtrlCmdHandler(signal []byte) error {
	return nil
}
