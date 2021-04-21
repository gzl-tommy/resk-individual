package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/consul"
	"github.com/tietang/props/v3/ini"
	"github.com/tietang/props/v3/kvs"
	_ "gzl-tommy/resk-individual"
	"gzl-tommy/resk-individual/infra"
	"gzl-tommy/resk-individual/infra/base"
	"net/http"
	_ "net/http/pprof"
)

//func main() {
//	//获取程序运行文件所在的路径
//	file := kvs.GetCurrentFilePath("config.ini", 1)
//	//加载和解析配置文件
//	conf := ini.NewIniFileCompositeConfigSource(file)
//	base.InitLog(conf)
//	app := infra.New(conf)
//	app.Start()
//}

func main() {

	// 通过 HTTP 服务来开启运行时性能剖析
	go func() {
		logrus.Info(http.ListenAndServe(":6060", nil))
	}()

	flag.Parse()
	profile := flag.Arg(0)
	if profile == "" {
		profile = "dev"
	}

	//获取程序运行文件所在的路径
	file := kvs.CurrentFilePath("boot.ini", 1)
	logrus.Info(file)

	//加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(file)
	if _, e := conf.Get("profile"); e != nil {
		conf.Set("profile", profile)
	}

	addr := conf.GetDefault("consul.address", "127.0.0.1:8500")
	contexts := conf.KeyValue("consul.contexts").Strings()
	logrus.Info("consul address:", addr)
	logrus.Info("consul contexts:", contexts)

	consulConf := consul.NewCompositeConsulConfigSourceByType(contexts, addr, kvs.ContentIni)
	consulConf.Add(conf)

	base.InitLog(consulConf)
	app := infra.New(consulConf)
	app.Start()
}
