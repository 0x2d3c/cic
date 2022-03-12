package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

const (
	Server = "s"
	Client = "c"
)

var (
	file, mode string
)

func Init() error {
	bs, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read config file err:%w", err)
	}

	if err = json.Unmarshal(bs, &cfg); err != nil {
		return fmt.Errorf("parse config file err:%w", err)
	}
	return nil
}

func main() {
	flag.StringVar(&mode, "m", "s", "running mode")
	flag.StringVar(&file, "f", "cfg.json", "used config file path")
	flag.Parse()

	if err := Init(); err != nil {
		panic(err)
	}

	switch mode {
	case Server:
		ServerRun()
	case Client:
		ClientRun()
	}
}
