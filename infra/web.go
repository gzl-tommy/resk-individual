package infra

var apiInitializerRegister *InitializeRegister = new(InitializeRegister)

// 获取注册的 web api 初始化对象
func GetApiInitializers() []Initializer {
	return apiInitializerRegister.Initializers
}

//注册WEB API初始化对象
func RegisterApi(ai Initializer) {
	apiInitializerRegister.Register(ai)
}

type WebApiStarter struct {
	BaseStarter
}

func (w *WebApiStarter) Setup(ctx StarterContext) {
	for _, v := range GetApiInitializers() {
		v.Init()
	}
}
