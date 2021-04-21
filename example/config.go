package main

import (
	"fmt"
	"github.com/tietang/props/v3/ini"
	"github.com/tietang/props/v3/kvs"
)

func main() {
	file := kvs.GetCurrentFilePath("config.ini", 1)
	conf := ini.NewIniFileConfigSource(file)
	port := conf.GetIntDefault("app.server.port", 18080)
	fmt.Println("====",port)
	fmt.Println(conf.GetDuration("app.time"))	
}
