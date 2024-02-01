package service

// maintainer 用于维护不同类型生命周期的服务实例创建及注入工作
type maintainer interface {

	// CreateServiceInstance 创建服务实例
	// 返回:
	// instance 服务实例
	// needInject 指示返回的instance实例是否需要后续执行注入操作
	CreateServiceInstance(provider *interimProvider) (instance any, needInject bool)

	// InjectForInstance 为服务实例执行注入操作
	InjectForInstance(v any, provider *interimProvider)

	// Fork 派生一个用于在下级Scope中使用的serviceMaintainer
	Fork(scope *defaultProvider) maintainer

	Clear()

	ReUse()
}

type maintainerBase struct {
	// 实例构建器
	Constructor Constructor

	// 实例注入器
	Injector Injector
}

func (m *maintainerBase) Clear() {
}

func (m *maintainerBase) ReUse() {
}
