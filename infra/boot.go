package infra

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
	"reflect"
)

// 应用程序的启动管理器
type BootApplication struct {
	IsTest     bool
	conf       kvs.ConfigSource
	starterCtx StarterContext
}

//构造系统
func New(conf kvs.ConfigSource) *BootApplication {
	b := &BootApplication{conf: conf, starterCtx: StarterContext{}}
	b.starterCtx.SetProps(conf)
	return b
}

func (b *BootApplication) Start() {
	//1. 初始化starter
	b.init()
	//2. 安装starter
	b.setup()
	//3. 启动starter
	b.start()
}

// 程序初始化
func (b *BootApplication) init() {
	log.Info("Initializing starters...")
	for _, v := range GetStarters() {
		log.Debugf("Initializing: PriorityGroup=%d,Priority=%d,type=%s", v.PriorityGroup(), v.Priority(), reflect.TypeOf(v).String())
		v.Init(b.starterCtx)
	}
}
func (b *BootApplication) setup() {
	log.Info("Setup starters...")
	for _, v := range GetStarters() {
		log.Debug("Setup: ", reflect.TypeOf(v).String())
		v.Setup(b.starterCtx)
	}
}
func (b *BootApplication) start() {
	log.Info("Starting starters...")
	for i, v := range GetStarters() {
		log.Debug("Starting: ", reflect.TypeOf(v).String())
		if v.StartBlocking() {
			if i+1 == len(GetStarters()) {
				v.Start(b.starterCtx)
			} else {
				go v.Start(b.starterCtx)
			}
		} else {
			v.Start(b.starterCtx)
		}
	}
}

//程序开始运行，开始接受调用
func (b *BootApplication) Stop() {
	log.Info("Stoping starters...")
	for _, v := range GetStarters() {
		log.Debug("Stoping: ", reflect.TypeOf(v).String())
		v.Stop(b.starterCtx)
	}
}
