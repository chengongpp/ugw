package ugwrt

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type RtInstance struct {
	Name           string
	WorkDir        string
	Args           []string
	MessageFormats map[string]*proto.Message
	LogLevel       log.Level
	Logger         []*log.Logger
	Outbounds      []Outbound
	Stats          *RtStat
	//TODO: Move to Inbound struct
	Host           string
	Port           int
	MaxConnections int
	Protocol       string
	Inbound        Inbound
}

func (rt *RtInstance) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	rt.Log(TxLog, log.InfoLevel, "Host=[%s] Protocol=[%s] Message=[%s]", request.Host, rt.Protocol, rt.Inbound.MessageFormat)
	switch rt.Protocol {
	case "httpJson":
		body, err := io.ReadAll(request.Body)
		if err != nil {

		}
		err = json.Unmarshal(body, make(interface{}, 1))
		if err != nil {

		}
		panic("")
	case "httpForm":
		panic("")
	case "httpXml":
		panic("")
	case "httpQuery":
		panic("")
	default:
		//TODO customized content protocol
	}
	//TODO implement me
	panic("implement me")
}

type Inbound struct {
	Name           string
	Host           string
	Port           int
	MaxConnections int
	MessageFormat  *proto.Message
}

type Outbound struct {
	Name          string
	Host          string
	MessageFormat *proto.Message
	Protocol      string
	Port          int
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

func NewRtInstance(conf Config) *RtInstance {
	var level log.Level
	switch strings.ToLower(conf.LogLevel) {
	case "debug":
		level = log.DebugLevel
	case "info":
		level = log.InfoLevel
	case "warn":
		level = log.WarnLevel
	case "error":
		level = log.ErrorLevel
	case "trace":
		level = log.FatalLevel
	default:
		fmt.Fprintf(os.Stderr, "Invalid log level: %s\n", conf.LogLevel)
	}
	log.SetLevel(level)

	loggers := make([]*log.Logger, 5)
	//Init AppLog
	appLogger := log.New()
	txLogger := log.New()
	detailLogger := log.New()
	traceLogger := log.New()
	appLogger.SetLevel(level)
	txLogger.SetLevel(level)
	detailLogger.SetLevel(level)
	traceLogger.SetLevel(level)
	loggers[0] = appLogger
	loggers[1] = txLogger
	loggers[2] = detailLogger
	loggers[3] = traceLogger
	switch conf.LogDir {
	case "":
		appLogger.SetOutput(os.Stdout)
		txLogger.SetOutput(os.Stdout)
		detailLogger.SetOutput(os.Stdout)
		traceLogger.SetOutput(os.Stdout)
	default:
		logPaths := []string{
			"app.log",
			"tx.log",
			"detail.log",
			"trace.log",
		}
		for i, filename := range logPaths {
			logFile, err := os.OpenFile(conf.LogDir+"/"+filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
				os.Exit(1)
			}
			loggers[i].SetOutput(logFile)
		}
	}
	return &RtInstance{
		Name:           conf.Name,
		WorkDir:        "",
		Args:           os.Args,
		Host:           conf.Host,
		Port:           conf.Port,
		Protocol:       conf.Protocol,
		MaxConnections: conf.MaxConnections,
		LogLevel:       level,
		Logger:         loggers,
		Outbounds:      conf.OutBounds,
	}
}

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
	if rt.Protocol == "http" {
		go rt.HttpServerLoop(mainCh)
	} else if rt.Protocol == "tcp" {
		go rt.TcpServerLoop(mainCh)
	} else {
		return fmt.Errorf("unsupported protocol: %s", rt.Protocol)
	}
	go rt.CtrlLoop(mainCh)

	<-mainCh
	return nil
}

func (rt *RtInstance) HttpServerLoop(mainCh chan error) {
	// TODO
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", rt.Host, rt.Port),
		Handler:           rt,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		MaxHeaderBytes:    0,
	}
	err := server.ListenAndServe()
	if err != nil {
		rt.Log(AppLog, log.ErrorLevel, "ListenHTTP failed: %v", err)
		mainCh <- err
		return
	}
}

func (rt *RtInstance) TcpServerLoop(mainCh chan error) {
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
