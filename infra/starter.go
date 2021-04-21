package infra

import (
	"github.com/tietang/props/v3/kvs"
	"reflect"
	"github.com/sirupsen/logrus"
	"sort"
)

const (
	KeyProps = "_conf"
)

// 资源启动器上下文，
// 用在服务资源初始化、安装、启动和停止的生命周期中变量和对象的传递
type StarterContext map[string]interface{}

func (s StarterContext) Props() kvs.ConfigSource {
	p := s[KeyProps]
	if p == nil {
		panic("配置还没被初始化")
	}
	return p.(kvs.ConfigSource)
}

func (s StarterContext) SetProps(conf kvs.ConfigSource) {
	s[KeyProps] = conf
}

//资源启动器，每个应用少不了依赖其他资源，比如数据库，缓存，消息中间件等等服务
//启动器实现类，不需要实现所有方法，只需要实现对应的阶段方法即可，可以嵌入@BaseStarter
//通过实现资源启动器接口和资源启动注册器，友好的管理这些资源的初始化、安装、启动和停止。
//Starter对象注册器，所有需要在系统启动时需要实例化和运行的逻辑，都可以实现此接口
//注意只有Start方法才能被阻塞，如果是阻塞Start()，同时StartBlocking()要返回true
type Starter interface {
	//资源初始化，通常把一些准备资源放在这里运行
	Init(StarterContext)

	//资源的安装，所有启动需要的具备条件，使得资源达到可以启动的就备状态
	Setup(StarterContext)

	//启动资源，达到可以使用的状态
	Start(StarterContext)

	//说明该资源启动器开始启动服务时，是否会阻塞
	//如果存在多个阻塞启动器时，只有最后一个阻塞，之前的会通过goroutine来异步启动
	//所以，需要规划好启动器注册顺序
	StartBlocking() bool

	//资源停止：
	//通常在启动时遇到异常时或者启用远程管理时，用于释放资源和终止资源的使用，
	//通常要优雅的释放，等待正在进行的任务继续，但不再接受新的任务
	Stop(StarterContext)

	PriorityGroup() PriorityGroup
	Priority() int
}

type PriorityGroup int

const (
	SystemGroup         PriorityGroup = 30
	BasicResourcesGroup PriorityGroup = 20
	AppGroup            PriorityGroup = 10

	INT_MAX          = int(^uint(0) >> 1)
	DEFAULT_PRIORITY = 10000
)

// 默认的空实现,方便资源启动器的实现
type BaseStarter struct {
}

func (b *BaseStarter) Init(ctx StarterContext)      {}
func (b *BaseStarter) Setup(ctx StarterContext)     {}
func (b *BaseStarter) Start(ctx StarterContext)     {}
func (b *BaseStarter) Stop(ctx StarterContext)      {}
func (b *BaseStarter) StartBlocking() bool          { return false }
func (s *BaseStarter) PriorityGroup() PriorityGroup { return BasicResourcesGroup }
func (s *BaseStarter) Priority() int                { return DEFAULT_PRIORITY }

//服务启动注册器
//不用需外部构造，全局只有一个
type starterRegister struct {
	nonBlockingStarters []Starter
	blockingStarters    []Starter
}

//启动器注册
func (r *starterRegister) Register(s Starter) {
	if s.StartBlocking() {
		r.blockingStarters = append(r.blockingStarters, s)
	} else {
		r.nonBlockingStarters = append(r.nonBlockingStarters, s)
	}

	logrus.Infof("Register starter: %s", reflect.TypeOf(s).String())
}

func (r *starterRegister) AllStarters() []Starter {
	starters := make([]Starter, 0)
	starters = append(starters, r.nonBlockingStarters...)
	starters = append(starters, r.blockingStarters...)
	return starters
}

var vStarterRegister *starterRegister = &starterRegister{}

// 注册starter
func Register(s Starter) {
	vStarterRegister.Register(s)
}

func SortStarters() {
	sort.Sort(Starters(vStarterRegister.AllStarters()))
}

//获取所有注册的starter
func GetStarters() []Starter {
	return vStarterRegister.AllStarters()
}

type Starters []Starter

func (s Starters) Len() int      { return len(s) }
func (s Starters) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Starters) Less(i, j int) bool {
	return s[i].PriorityGroup() > s[j].PriorityGroup() && s[i].Priority() > s[j].Priority()
}
