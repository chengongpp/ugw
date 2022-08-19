package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	toml "github.com/BurntSushi/toml"
	"github.com/chengongpp/ugw/pkg/ugwrt"
	"os"
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
			_, _ = fmt.Fprintf(os.Stderr, "create config file failed: %v\n", err)
			os.Exit(1)
		}
		defer func(confFile *os.File) {
			_ = confFile.Close()
		}(confFile)
		if _, err := confFile.WriteString(configTmpl); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "write config file failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	var conf ugwrt.Config = ugwrt.Config{
		LogDir: "",
	}

	if _, err := os.Stat(configPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_, _ = fmt.Fprintf(os.Stderr, "Config file %s not found. -h For help\n", configPath)
			os.Exit(1)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Config file %s open failed.\n", configPath)
			os.Exit(1)
		}
	}

	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	rt := ugwrt.NewRtInstance(conf)
	err := rt.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error running rt: %v\n", err)
		os.Exit(1)
	}
}
