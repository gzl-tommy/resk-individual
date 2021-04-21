package lb

import (
	"fmt"
	"github.com/gzl-tommy/go-eureka-client/eureka"
	"strings"
)

type Apps struct {
	Client *eureka.Client
}

func (as *Apps) Get(appName string) *App {
	var eApp *eureka.Application
	for _, a := range as.Client.Applications.Applications {
		if a.Name == strings.ToUpper(appName) {
			eApp = &a
		}
	}
	if eApp == nil {
		return nil
	}

	app := &App{
		Name:      eApp.Name,
		Instances: make([]*ServerInstance, 0),
		lb:        &RoundRobinBalancer{},
	}
	for _, ins := range eApp.Instances {
		var port int
		if ins.SecurePort.Enabled {
			port = ins.SecurePort.Port
		} else {
			port = ins.Port.Port
		}
		si := &ServerInstance{
			InstanceId: ins.InstanceId,
			AppName:    appName,
			Address:    fmt.Sprintf("%s:%d", ins.IpAddr, port),
			Status:     Status(ins.Status),
		}
		app.Instances = append(app.Instances, si)
	}
	return app
}

type App struct {
	Name      string
	Instances []*ServerInstance
	lb        Balancer
}

func (a *App) Get(key string) *ServerInstance {
	ins := a.lb.Next(key, a.Instances)
	return ins
}

// 服务实例状态
type Status = string

const (
	StatusEnabled  Status = "enabled"
	StatusDisabled Status = "disabled"
)

// 服务实例
type ServerInstance struct {
	InstanceId string
	AppName    string
	Address    string
	Status     Status
}
