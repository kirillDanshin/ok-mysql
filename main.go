package main

import (
	"flag"

	"github.com/kirillDanshin/myutils"
	"github.com/kirillDanshin/ok-mysql/ok"
	"github.com/spf13/viper"
)

var (
	cfg = viper.New()

	configDir        = flag.String("configDir", ".", "path to yaml config file")
	configName       = flag.String("configName", "config", "yaml config file without extension")
	initConfigNeeded = flag.Bool("init", false, "init config file")
)

func main() {
	// flag.Parse()
	// cfg.AddConfigPath(*configDir)
	// cfg.SetConfigName(*configName)
	// cfg.SetConfigType("yaml")
	// cfg.ReadInConfig()
	//
	// if *initConfigNeeded {
	// 	// initConfig(*configDir, *configName) //PLANNING<kirillDanshin>
	// 	return
	// }
	// dLog("cfg.AllSettings(): %#+v\n", cfg.AllSettings())

	flag.Parse()

	instance, err := ok.NewInstance(
		&ok.Config{
			Address: "127.0.0.1:3306",
		},
	)
	myutils.LogFatalError(err)
	err = instance.Run()
	myutils.LogFatalError(err)

}
