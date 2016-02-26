package main

import (
	"flag"
	"fmt"
)

type Arg struct {
	Endpoint     string
	Mode         string
	ProxyPort    string
	StubFileName string
}

func parseArg() {
	flag.StringVar(&arg.Endpoint, "target", "https://gliese1dtac-blqrysb.tac.co.th:7834/", "target URL for reverse proxy")
	flag.StringVar(&arg.ProxyPort, "port", "9000", "proxy running on port. EX: 9000")
	flag.StringVar(&arg.Mode, "mode", "Record", "proxy running mode [Record/Replay], default is Replay")
	flag.StringVar(&arg.StubFileName, "stubFileName", "stub.json", "record to file name EX stub.json")
	flag.Parse()

	println("===========================================")
	println(arg.Mode, "Mode")
	println("===========================================")
	fmt.Printf("arg = %#v\n", arg)
}
