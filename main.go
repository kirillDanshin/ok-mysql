package main

import (
	"flag"

	"github.com/kirillDanshin/myutils"
	"github.com/kirillDanshin/ok-mysql/ok"
)

var (
	addr = flag.String("addr", "", "address to listen (required)")
)

func main() {
	flag.Parse()
	myutils.RequiredStrFatal("address", *addr)

	instance, err := ok.NewInstance(
		&ok.Config{
			Address: *addr,
		},
	)
	myutils.LogFatalError(err)
	err = instance.Run()
	myutils.LogFatalError(err)

}
