package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	toml "github.com/BurntSushi/toml"
	"github.com/chengongpp/ugw/pkg/ugwrt"
	log "github.com/sirupsen/logrus"
)

var version = "UGW 0.0.1 by chengongpp"
var configPath string
var showVersion bool
var generateConfig string

//go:embed config.toml
var configTmpl string

func init() {
	flag.StringVar(&configPath, "c", "config.toml", "config file path")
	flag.StringVar(&configPath, "config", "config.toml", "config file path")
	flag.BoolVar(&showVersion, "V", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.StringVar(&generateConfig, "n", "", "generate config file")
}

func main() {
	flag.Parse()
	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	if generateConfig != "" {
		confFile, err := os.Create(generateConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create config file failed: %v\n", err)
			os.Exit(1)
		}
		defer confFile.Close()
		if _, err := confFile.WriteString(configTmpl); err != nil {
			fmt.Fprintf(os.Stderr, "write config file failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	var conf ugwrt.Config = ugwrt.Config{
		LogDir: "",
	}

	if _, err := os.Stat(configPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "Config file %s not found. -h For help\n", configPath)
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "Config file %s open failed.\n", configPath)
			os.Exit(1)
		}
	}

	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
		os.Exit(1)
	}
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
				fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
				os.Exit(1)
			}
			loggers[i].SetOutput(logFile)
		}
	}
	runtime := ugwrt.RtInstance{
		Name:           conf.Name,
		WorkDir:        "",
		Args:           os.Args,
		Host:           conf.Host,
		Port:           conf.Port,
		MaxConnections: conf.MaxConnections,
		LogLevel:       level,
		Logger:         loggers,
		OutBounds:      conf.OutBounds,
	}
	runtime.Run()
}
